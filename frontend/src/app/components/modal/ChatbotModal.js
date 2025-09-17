import { useState, useRef, useEffect } from "react";
import Modal from "./Modal";
import Button from "../common/Button";
import Text from "../common/Text";
import TextInput from "../Input/TextInput";
import TextareaInput from "../Input/TextAreaInput";
import { FaRobot, FaPaperPlane, FaSpinner, FaUser } from "react-icons/fa";
import Subtitle from "../common/Subtitle";

function ChatMessage({ role, content }) {
  return (
    <div
      className={`flex items-end gap-2 max-w-[80%] mb-2 px-4 py-2 rounded-lg ${
        role === "user"
          ? "bg-blue-100 ml-auto flex-row-reverse text-right"
          : "bg-green-100 mr-auto text-left"
      }`}
    >
      {role === "user" && <FaUser className={"text-blue-500"} size={20} />}
      <Text
        textSize="md"
        textColor={role === "user" ? "text-blue-600" : "text-green-600"}
        className="ml-2"
      >
        {content}
      </Text>
    </div>
  );
}

function ChatInput({ value, onChange, onSend, loading }) {
  const handleKeyDown = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (!loading && value) {
        onSend();
      }
    }
  };

  return (
    <div className="flex gap-2 items-end">
      <TextareaInput
        name="ai-chat-input"
        placeholder="Unesite pitanje..."
        value={value}
        onChange={onChange}
        className="flex-1"
        disabled={loading}
        rows={2}
        onKeyDown={handleKeyDown}
      />
      <Button
        onClick={onSend}
        disabled={loading || !value}
        icon={loading ? FaSpinner : FaPaperPlane}
        iconClassName={loading ? "animate-spin" : ""}
        color="secondary"
        className="h-10"
      >
        {loading ? "Generisanje..." : "Pošalji"}
      </Button>
    </div>
  );
}

// AIChatModal component
export default function AIChatModal({ open, onClose, accessToken }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const sessionIdRef = useRef(null);

  const sendMessage = async () => {
    if (!input) return;
    setLoading(true);
    const payload = { question: input };
    if (sessionIdRef.current) payload.session_id = sessionIdRef.current;
    const resp = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/chat`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(payload),
      },
    );
    const data = await resp.json();
    if (!sessionIdRef.current) sessionIdRef.current = data.session_id;
    setMessages((prev) => [
      ...prev,
      { role: "user", content: input },
      { role: "ai", content: data.answer },
    ]);
    setInput("");
    setLoading(false);
  };

  useEffect(() => {
    if (open) {
      sessionIdRef.current = null;
      setMessages([]);
    }
  }, [open]);

  return (
    <Modal open={open} onClose={onClose}>
      <Subtitle icon={FaRobot} showLine={false} className="mb-8">
        eDnevnik Asistent
      </Subtitle>
      <div className="max-h-64 overflow-y-auto mb-2 mt-4">
        {messages.map((msg, idx) => (
          <ChatMessage key={idx} role={msg.role} content={msg.content} />
        ))}
      </div>
      <ChatInput
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onSend={sendMessage}
        loading={loading}
      />
    </Modal>
  );
}
