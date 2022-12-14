import React, { useCallback, useEffect, useRef } from "react";

import autoAnimate from "@formkit/auto-animate";

import MessageComponent, { MessageProps } from "@components/Message";
import MessageSceleton from "@components/Sceletons/MessageSceleton";

import { MdMessage } from "react-icons/md";

type ChatProps = {
  messages: MessageProps[] | undefined;
  loading?: boolean;
};
//TODO: refactor messages to store them in state and check if they are on
//db instead of taking them from db and posting only after they are created with backend
function Chat({ messages, loading }: ChatProps) {
  const chatRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (chatRef.current) {
      chatRef.current?.scrollTo(0, chatRef.current.scrollHeight);
      autoAnimate(chatRef.current);
    }
  }, []);
  //FIXME: new chat doesn't display first message
  const renderMessages = () => {
    if (loading) {
      return [...Array(10)].map((_, i) => <MessageSceleton key={i} />);
    }
    if (!messages || messages.length === 0) {
      return (
        <div className="text-bold flex self-center text-lg ">
          <span>There are no messages </span>{" "}
          <MdMessage className="mx-2 self-center" />
          <span>in chat</span>
        </div>
      );
    }

    return messages?.map(({ body, createdAt, user, sent }) => (
      <MessageComponent
        body={body}
        user={user}
        createdAt={createdAt}
        sent={sent}
        key={body + createdAt.getTime()}
      />
    ));
  };

  return (
    <div
      ref={chatRef}
      className={`${
        loading || !messages ? "overflow-y-hidden" : "overflow-y-scroll"
      } custom-scroll flex w-full flex-grow flex-col-reverse   p-4 md:py-8 md:px-16`}
    >
      {renderMessages()}
    </div>
  );
}

export default Chat;
