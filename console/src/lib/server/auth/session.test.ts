import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  generateSessionId,
  createSession,
  validateSession,
  revokeSession,
  revokeAllSessionsForUser,
} from "./session";
import { db } from "$lib/server/db";
import { sessions, users } from "$lib/server/db/schema";
import * as drizzleOrm from "drizzle-orm";

const mockInsertValues = vi.fn();
const mockSelectLimit = vi.fn();
const mockSelectInnerJoin = vi.fn().mockReturnValue({ limit: mockSelectLimit });
const mockSelectWhere = vi
  .fn()
  .mockReturnValue({ innerJoin: mockSelectInnerJoin });
const mockSelectFrom = vi.fn().mockReturnValue({ where: mockSelectWhere });

const mockUpdateWhere = vi.fn();
const mockUpdateSet = vi.fn().mockReturnValue({ where: mockUpdateWhere });

vi.mock("$lib/server/db", () => ({
  db: {
    insert: vi.fn(() => ({ values: mockInsertValues })),
    select: vi.fn(() => ({ from: mockSelectFrom })),
    update: vi.fn(() => ({ set: mockUpdateSet })),
  },
}));

vi.mock("$lib/server/db/schema", () => ({
  sessions: {
    id: "sessions.id",
    userId: "sessions.userId",
    expiresAt: "sessions.expiresAt",
    revokedAt: "sessions.revokedAt",
  },
  users: {
    id: "users.id",
  },
}));

vi.mock("drizzle-orm", () => ({
  eq: vi.fn((a, b) => ({ op: "eq", args: [a, b] })),
  and: vi.fn((...args) => ({ op: "and", args })),
  isNull: vi.fn((a) => ({ op: "isNull", args: [a] })),
  gt: vi.fn((a, b) => ({ op: "gt", args: [a, b] })),
}));

describe("session.ts", () => {
  beforeEach(() => {
    vi.clearAllMocks();

    // If crypto is not available in test env, mock it
    if (typeof globalThis.crypto === "undefined") {
      globalThis.crypto = {
        getRandomValues: (arr: Uint8Array) => {
          for (let i = 0; i < arr.length; i++) {
            arr[i] = Math.floor(Math.random() * 256);
          }
          return arr;
        },
      } as any;
    }
  });

  describe("generateSessionId", () => {
    it("should generate a 64-character hex string", () => {
      const id = generateSessionId();
      expect(typeof id).toBe("string");
      expect(id).toHaveLength(64);
      expect(/^[0-9a-f]{64}$/.test(id)).toBe(true);
    });

    it("should generate unique ids", () => {
      const id1 = generateSessionId();
      const id2 = generateSessionId();
      expect(id1).not.toBe(id2);
    });
  });

  describe("createSession", () => {
    it("should create a session and insert it into the database", async () => {
      const result = await createSession("user_123");

      expect(result).toHaveProperty("id");
      expect(result).toHaveProperty("expiresAt");
      expect(result.id).toHaveLength(64);
      expect(result.expiresAt).toBeInstanceOf(Date);

      expect(db.insert).toHaveBeenCalledWith(sessions);
      expect(mockInsertValues).toHaveBeenCalledWith({
        id: result.id,
        userId: "user_123",
        expiresAt: result.expiresAt,
      });
    });
  });

  describe("validateSession", () => {
    it("should return session data if valid", async () => {
      const mockSessionData = {
        session: { id: "sess_1" },
        user: { id: "user_1" },
      };
      mockSelectLimit.mockResolvedValueOnce([mockSessionData]);

      const result = await validateSession("sess_1");

      expect(db.select).toHaveBeenCalled();
      expect(mockSelectFrom).toHaveBeenCalledWith(sessions);
      expect(mockSelectWhere).toHaveBeenCalled();
      expect(mockSelectInnerJoin).toHaveBeenCalledWith(
        users,
        expect.any(Object),
      );
      expect(mockSelectLimit).toHaveBeenCalledWith(1);

      expect(result).toEqual(mockSessionData);
    });

    it("should return null if session is invalid", async () => {
      mockSelectLimit.mockResolvedValueOnce([]);

      const result = await validateSession("invalid_sess");

      expect(result).toBeNull();
    });
  });

  describe("revokeSession", () => {
    it("should update revokedAt for the specific session", async () => {
      await revokeSession("sess_123");

      expect(db.update).toHaveBeenCalledWith(sessions);
      expect(mockUpdateSet).toHaveBeenCalledWith({
        revokedAt: expect.any(Date),
      });
      expect(mockUpdateWhere).toHaveBeenCalledWith(
        expect.objectContaining({
          op: "eq",
          args: ["sessions.id", "sess_123"],
        }),
      );
    });
  });

  describe("revokeAllSessionsForUser", () => {
    it("should update revokedAt for all active sessions of a user", async () => {
      await revokeAllSessionsForUser("user_456");

      expect(db.update).toHaveBeenCalledWith(sessions);
      expect(mockUpdateSet).toHaveBeenCalledWith({
        revokedAt: expect.any(Date),
      });
      expect(mockUpdateWhere).toHaveBeenCalledWith(
        expect.objectContaining({ op: "and" }),
      );
    });
  });
});
