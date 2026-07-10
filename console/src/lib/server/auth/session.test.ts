import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  createSession,
  validateSession,
  revokeSession,
  revokeAllSessionsForUser,
} from "./session";
import { db } from "$lib/server/db";
import { sessions, users } from "$lib/server/db/schema";
import * as drizzleOrm from "drizzle-orm";
import * as tokenModule from "./token";

const mockInsertValues = vi.fn();
const mockSelectLimit = vi.fn();
const mockSelectInnerJoin = vi.fn().mockReturnValue({ limit: mockSelectLimit });
const mockSelectWhere = vi
  .fn()
  .mockReturnValue({ innerJoin: mockSelectInnerJoin });
const mockSelectFrom = vi.fn().mockReturnValue({ where: mockSelectWhere });

const mockDeleteWhere = vi.fn();
const mockDelete = vi.fn().mockReturnValue({ where: mockDeleteWhere });

vi.mock("$lib/server/db", () => ({
  db: {
    insert: vi.fn(() => ({ values: mockInsertValues })),
    select: vi.fn(() => ({ from: mockSelectFrom })),
    delete: mockDelete,
  },
}));

vi.mock("$lib/server/db/schema", () => ({
  sessions: {
    id: "sessions.id",
    userId: "sessions.userId",
    expiresAt: "sessions.expiresAt",
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
  });

  describe("createSession", () => {
    it("should create a session and insert it into the database", async () => {
      const result = await createSession("user_123");

      expect(result).toHaveProperty("token");
      expect(result).toHaveProperty("expiresAt");
      expect(typeof result.token).toBe("string");
      expect(result.expiresAt).toBeInstanceOf(Date);

      expect(db.insert).toHaveBeenCalledWith(sessions);
      expect(mockInsertValues).toHaveBeenCalledWith({
        id: expect.any(String),
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
    it("should delete the specific session", async () => {
      await revokeSession("sess_123");

      expect(db.delete).toHaveBeenCalledWith(sessions);
      expect(mockDeleteWhere).toHaveBeenCalledWith(
        expect.objectContaining({
          op: "eq",
          args: ["sessions.id", expect.any(String)],
        }),
      );
    });
  });

  describe("revokeAllSessionsForUser", () => {
    it("should delete all active sessions of a user", async () => {
      await revokeAllSessionsForUser("user_456");

      expect(db.delete).toHaveBeenCalledWith(sessions);
      expect(mockDeleteWhere).toHaveBeenCalledWith(
        expect.objectContaining({
          op: "eq",
          args: ["sessions.userId", "user_456"],
        }),
      );
    });
  });
});
