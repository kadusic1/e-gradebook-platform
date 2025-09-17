"use client";
import { useState } from "react";
import ChatbotModal from "./ChatbotModal";
import Button from "../common/Button";
import { FaRobot } from "react-icons/fa";

export default function AIChatEntry({ accessToken, open, setOpen }) {
  return (
    <>
      {!open ? (
        <Button
          className="fixed bottom-8 right-6 z-[50]"
          icon={FaRobot}
          color="secondary"
          onClick={() => setOpen(true)}
        >
          Asistent
        </Button>
      ) : (
        <ChatbotModal
          open={open}
          onClose={() => setOpen(false)}
          accessToken={accessToken}
        />
      )}
    </>
  );
}
