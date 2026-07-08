CREATE TABLE IF NOT EXISTS console.invitations (
    id TEXT PRIMARY KEY,
    integrator_id UUID NOT NULL REFERENCES public.api_integrators(id),
    invited_by UUID NOT NULL REFERENCES console.users(id),
    email TEXT NOT NULL,
    role console.user_role NOT NULL DEFAULT 'member',
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
