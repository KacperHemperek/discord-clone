import { Context } from "@server/trpc/context";

import { prisma } from "@server/db/client";

export const getWelcomeChannel = async (): Promise<number> => {
  try {
    const channel = await prisma?.channel.findUnique({
      where: { name: "Welcome" },
    });
    if (!channel?.id) {
      throw new Error("There is no Welcome channel please create one");
    }
    return channel.id;
  } catch (err) {
    console.error(err);
    throw new Error("Couldn't retrive Welcome channel id");
  }
};
