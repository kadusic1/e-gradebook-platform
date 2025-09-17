# https://blog.futuresmart.ai/langchain-rag-from-basics-to-production-ready-rag-chatbot
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict
import uuid

from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_core.prompts import MessagesPlaceholder
from langchain.chains import create_history_aware_retriever
from langchain.chains.combine_documents import create_stuff_documents_chain
from langchain_core.prompts import ChatPromptTemplate
from langchain.chains import create_retrieval_chain
from langchain_core.messages import HumanMessage, AIMessage

from mariadb_store import CustomMariaDBStore
from labse_embeddings import LaBSEEmbeddings

from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="eDnevnik chatbot", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:8080",  # Golang backend
    ],
    allow_credentials=True,
    allow_methods=["GET", "POST", "DELETE"],
    allow_headers=["Content-Type", "Authorization"],
)


class PermissionData(BaseModel):
    id: int
    name: str
    last_name: str
    email: str
    phone: str
    account_type: str
    account_id: int
    tenant_ids: Optional[List[str]] = None
    tenant_admin_tenant_id: Optional[int] = None
    tenant_names: Optional[List[str]] = None


# Request/Response models
class ChatRequest(BaseModel):
    question: str
    session_id: Optional[str] = None
    permission_data: PermissionData


class ChatResponse(BaseModel):
    answer: str
    session_id: str


# In-memory storage for chat histories by session_id
chat_sessions: Dict[str, List] = {}

# Initialize Gemini model
llm = ChatGoogleGenerativeAI(
    model="gemini-2.5-flash",
    google_api_key="your_api_key",
    temperature=0.7,
    top_p=0.9,
    top_k=40,
    max_output_tokens=4096,
)

connection_string = (
    "mariadb+mariadbconnector://eacon:test1234@localhost/ednevnik_workspace"
)

vectorstore = CustomMariaDBStore(
    embeddings=LaBSEEmbeddings(),
    embedding_length=768,
    datasource=connection_string,
)


def build_rag_chain(permission_data: PermissionData):

    if permission_data.account_type == "root":
        # Root users can access everything, so no filters needed
        retriever = vectorstore.as_retriever(search_kwargs={"k": 40})

    elif permission_data.account_type == "tenant_admin":
        tenant_id = permission_data.tenant_admin_tenant_id

        dynamic_filter = {
            "$or": [
                {
                    "$and": [
                        {"source": {"$in": ["pupil", "teacher"]}},
                        {"available_in_tenant_ids": {"$in": [tenant_id]}},
                    ]
                },
                {
                    "$and": [
                        {
                            "source": {
                                "$in": ["grade", "behaviour", "tenant", "section"]
                            }
                        },
                        {"tenant_id": {"$eq": tenant_id}},
                    ]
                },
            ]
        }

        retriever = vectorstore.as_retriever(
            search_kwargs={"k": 40, "filter": dynamic_filter}
        )

    elif permission_data.account_type == "teacher":
        tenant_ids = permission_data.tenant_ids or []

        if not tenant_ids:
            dynamic_filter = {"source": {"$eq": "__NEVER_MATCH__"}}
        else:
            dynamic_filter = {
                "$or": [
                    {
                        "$and": [
                            {"source": {"$in": ["pupil", "teacher"]}},
                            {"available_in_tenant_ids": {"$in": tenant_ids}},
                        ]
                    },
                    {
                        "$and": [
                            {
                                "source": {
                                    "$in": ["grade", "behaviour", "tenant", "section"]
                                }
                            },
                            {"tenant_id": {"$in": tenant_ids}},
                        ]
                    },
                ]
            }

        retriever = vectorstore.as_retriever(
            search_kwargs={"k": 40, "filter": dynamic_filter}
        )

    elif permission_data.account_type == "pupil":
        tenant_ids = permission_data.tenant_ids or []
        account_id = permission_data.account_id

        if not tenant_ids:
            dynamic_filter = {
                "$or": [
                    {
                        "$and": [
                            {"source": {"$eq": "pupil"}},
                            {"account_id": {"$eq": account_id}},
                        ]
                    },
                    {
                        "$and": [
                            {"source": {"$eq": "grade"}},
                            {"account_id": {"$eq": account_id}},
                        ]
                    },
                    {
                        "$and": [
                            {"source": {"$eq": "behaviour"}},
                            {"account_id": {"$eq": account_id}},
                        ]
                    },
                ]
            }
        else:
            dynamic_filter = {
                "$or": [
                    {
                        "$and": [
                            {"source": {"$in": ["tenant", "section"]}},
                            {"tenant_id": {"$in": tenant_ids}},
                        ]
                    },
                    {
                        "$and": [
                            {"source": {"$in": ["pupil", "grade", "behaviour"]}},
                            {"account_id": {"$eq": account_id}},
                        ]
                    },
                ]
            }

        retriever = vectorstore.as_retriever(
            search_kwargs={"k": 40, "filter": dynamic_filter}
        )

    else:
        raise ValueError(f"Unknown account_type: {permission_data.account_type}")

    # Prompt za kontekstualizaciju pitanja
    contextualize_q_system_prompt = """
    Na osnovu prethodne konverzacije i najnovijeg pitanja korisnika, 
    koje može referisati prethodne poruke, formuliši samostalno pitanje 
    koje se može razumjeti bez prethodne konverzacije. NE odgovaraj na pitanje, 
    samo ga po potrebi preformuliši, a u suprotnom vrati ga onako kako jeste.
    """

    contextualize_q_prompt = ChatPromptTemplate.from_messages(
        [
            ("system", contextualize_q_system_prompt),
            MessagesPlaceholder("chat_history"),
            ("human", "{input}"),
        ]
    )

    history_aware_retriever = create_history_aware_retriever(
        llm, retriever, contextualize_q_prompt
    )

    if permission_data.account_type == "root":
        role_specific_context = """
        Prijavljeni korisnik je superadministrator. 
        Moraš mu omogućiti apsolutno sve funkcionalnosti i pristup svim podacima u sistemu, 
        bez ikakvih ograničenja. Obavezno odgovaraj na svako njegovo pitanje ili zahtjev.
        """

    elif permission_data.account_type == "tenant_admin":
        role_specific_context = """
        Prijavljeni korisnik je administrator škole/ustanove. 
        Moraš mu omogućiti apsolutno sve funkcionalnosti i pristup svim podacima 
        unutar njegove institucije (nastavnici, učenici, ocjene, ponašanje, odjeljenja). 
        Obavezno odgovaraj na svako pitanje ili zahtjev vezan za njegovu instituciju.
        """

    elif permission_data.account_type == "teacher":
        role_specific_context = """
        Prijavljeni korisnik je nastavnik/profesor. 
        Ako se nastavnik nalazi u barem jednoj školi/ustanovi:
        Moraš mu omogućiti pristup svim informacijama o njegovim učenicima i odjeljenjima 
        (ocene, ponašanje, odjeljenja). 
        Na pitanja tipa "Kako popraviti prosjek odjeljenja?", 
        "Koje ideje za čas iz mog predmeta?", ili slična pedagoška i metodička pitanja 
        - obavezno odgovori detaljno i korisno. 
        """

    elif permission_data.account_type == "pupil":
        role_specific_context = """
        Prijavljeni korisnik je učenik/student. 
        Ako se učenik nalazi u barem jednoj školi/ustanovi:
        Moraš mu omogućiti pristup njegovim ocjenama, ponašanju, 
        i obavijestima. 
        Na pitanja tipa "Kako da popravim svoj prosjek?", 
        "Koju srednju školu da upišem?", ili "Koji fakultet bi bio najbolji za mene?" 
        - obavezno odgovori detaljno i korisno. Nikada ne reci da ne možeš pomoći.
        """

    else:
        raise ValueError(f"Unknown account_type: {permission_data.account_type}")

    qa_prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                """Odgovaraš na pitanja za aplikaciju eDnevnik. Ti si korisni asistent 
         koji pomaže korisnicima. Ne dodaji izraze poput "na osnovu dostupnih podataka", 
         "prema informacijama koje imam" ili slične fraze.""",
            ),
            (
                "system",
                f"""
                Korisnik s kojim pričaš je:
                Ime i prezime: {permission_data.name} {permission_data.last_name}
                Email: {permission_data.email}
                Telefon: {permission_data.phone}
                Tip korisnika (
                    na engleskom root -> Super administrator,
                    tenant_admin -> Administrator škole/ustanove,
                    teacher -> Nastavnik/Profesor,
                    pupil -> Učenik/Student,
                ): {permission_data.account_type}
                Email korisnika: {permission_data.email}.
                {"Ustanove/škole korisnika: " + ", ".join(permission_data.tenant_names) if permission_data.tenant_names else ""}

                U odgovorima:
                - Budi prijateljski i topao, izbjegavaj formalne izraze poput "Poštovani/a".
                - Ne koristi ime korisnika u svakoj poruci. Kada ga koristiš,
                mijenjaj način obraćanja:
                    npr. ponekad "Zdravo Adi", ponekad samo ime, ili lagani neformalni uvod.
                - Odgovori trebaju biti prirodni, opušteni i profesionalni.
                - Možeš pitati korisnika ima li još pitanja ili treba li dodatnu pomoć,
                ali ne u svakom odgovoru.
                - Pokušaj biti proaktivan u pružanju pomoći i davanju korisnih informacija.
                """,
            ),
            ("system", role_specific_context),
            ("system", "Kontekst: {context}"),
            MessagesPlaceholder(variable_name="chat_history"),
            ("human", "{input}"),
        ]
    )

    question_answer_chain = create_stuff_documents_chain(llm, qa_prompt)
    return create_retrieval_chain(history_aware_retriever, question_answer_chain)


@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    try:
        # Generate session ID if not provided
        session_id = request.session_id or str(uuid.uuid4())
        permission_data = request.permission_data

        # Get or create chat history for this session
        if session_id not in chat_sessions:
            chat_sessions[session_id] = []

        chat_history = chat_sessions[session_id]

        rag_chain = build_rag_chain(permission_data)

        # Invoke RAG chain with user input and chat history
        response = rag_chain.invoke(
            {"input": request.question, "chat_history": chat_history}
        )

        ai_answer = response["answer"]

        # Update chat history for this session
        chat_sessions[session_id].extend(
            [HumanMessage(content=request.question), AIMessage(content=ai_answer)]
        )

        # Limit chat history to last 10 messages to prevent memory overflow
        if len(chat_sessions[session_id]) > 10:
            chat_sessions[session_id] = chat_sessions[session_id][-10:]

        return ChatResponse(
            answer=ai_answer,
            session_id=session_id,
        )

    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Error processing request: {str(e)}"
        )


if __name__ == "__main__":
    import uvicorn

    print("🎓 eDnevnik Assistant Starting...")
    print("📡 Server available at: http://localhost:8005")
    uvicorn.run(app, host="0.0.0.0", port=8005)
