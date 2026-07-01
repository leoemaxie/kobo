import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebar: SidebarsConfig = {
  apisidebar: [
    {
      type: "doc",
      id: "api/kobo-api",
    },
    {
      type: "category",
      label: "Identities",
      link: {
        type: "doc",
        id: "api/identities",
      },
      items: [
        {
          type: "doc",
          id: "api/create-identity",
          label: "Register a new identity and provision its virtual account.",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api/get-identity",
          label: "Fetch an identity record, its current state, and linked account.",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api/update-identity",
          label: "Update display profile fields (rename, contact details).",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "api/close-identity",
          label: "Initiate closure of an identity's virtual account.",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api/reopen-identity",
          label: "Reopen a CLOSED identity's account.",
          className: "api-method post",
        },
      ],
    },
    {
      type: "category",
      label: "Accounts",
      link: {
        type: "doc",
        id: "api/accounts",
      },
      items: [
        {
          type: "doc",
          id: "api/list-account-transactions",
          label: "Paginated, reconciled transaction history for an account.",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api/get-account-statement",
          label: "Structured statement for a given period.",
          className: "api-method get",
        },
      ],
    },
    {
      type: "category",
      label: "Exceptions",
      link: {
        type: "doc",
        id: "api/exceptions",
      },
      items: [
        {
          type: "doc",
          id: "api/list-exceptions",
          label: "List unresolved misdirected-payment or unmatched-transfer cases.",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api/resolve-exception",
          label: "Apply a resolution to a flagged exception.",
          className: "api-method post",
        },
      ],
    },
    {
      type: "category",
      label: "Webhooks",
      link: {
        type: "doc",
        id: "api/webhooks",
      },
      items: [
        {
          type: "doc",
          id: "api/receive-nomba-webhook",
          label: "Internal endpoint Nomba calls into. Not used by Kobo integrators.",
          className: "api-method post",
        },
      ],
    },
    {
      type: "category",
      label: "System",
      link: {
        type: "doc",
        id: "api/system",
      },
      items: [
        {
          type: "doc",
          id: "api/health-check",
          label: "Liveness/readiness check.",
          className: "api-method get",
        },
      ],
    },
  ],
};

export default sidebar.apisidebar;
