
CREATE TABLE public.contas (
  id bigint PRIMARY KEY NOT NULL,
  cliente_id bigint NOT NULL,
  total bigint NOT NULL,
  limite bigint NOT NULL
);

CREATE TABLE public.transacoes (
  id bigint PRIMARY KEY NOT NULL,
  conta_id bigint NOT NULL,
  valor bigint NOT NULL,
  tipo character varying(1) NOT NULL,
  descricao character varying(255) NOT NULL,
  realizado_em timestamp without time zone NOT NULL
);

CREATE SEQUENCE public.transacoes_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

ALTER SEQUENCE public.transacoes_id_seq OWNED BY public.transacoes.id;

ALTER TABLE ONLY public.transacoes ALTER COLUMN id SET DEFAULT nextval('public.transacoes_id_seq'::regclass);

ALTER TABLE ONLY public.transacoes
  ADD CONSTRAINT fk_transacoes_conta_id FOREIGN KEY (conta_id) REFERENCES public.contas(id);

INSERT INTO public.contas (id, cliente_id, total, limite) VALUES (1, 1, 0, 100000);
INSERT INTO public.contas (id, cliente_id, total, limite) VALUES (2, 2, 0, 80000);
INSERT INTO public.contas (id, cliente_id, total, limite) VALUES (3, 3, 0, 1000000);
INSERT INTO public.contas (id, cliente_id, total, limite) VALUES (4, 4, 0, 10000000);
INSERT INTO public.contas (id, cliente_id, total, limite) VALUES (5, 5, 0, 500000);
