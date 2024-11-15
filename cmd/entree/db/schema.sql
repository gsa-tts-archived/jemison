SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: scheme; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.scheme AS ENUM (
    'http',
    'https'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: content_types; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.content_types (
    id integer NOT NULL,
    content_type text
);


--
-- Name: content_types_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.content_types ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.content_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: guestbook; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.guestbook (
    id bigint NOT NULL,
    scheme public.scheme NOT NULL,
    host bigint,
    path text NOT NULL,
    content_sha1 text,
    content_length integer,
    content_type integer,
    last_updated date,
    last_fetched date,
    next_fetch date
);


--
-- Name: guestbook_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.guestbook ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.guestbook_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: hosts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hosts (
    id bigint NOT NULL,
    host text
);


--
-- Name: hosts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.hosts ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.hosts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: content_types content_types_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.content_types
    ADD CONSTRAINT content_types_pkey PRIMARY KEY (id);


--
-- Name: guestbook guestbook_host_path_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guestbook
    ADD CONSTRAINT guestbook_host_path_key UNIQUE (host, path);


--
-- Name: guestbook guestbook_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guestbook
    ADD CONSTRAINT guestbook_pkey PRIMARY KEY (id);


--
-- Name: hosts hosts_host_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_host_key UNIQUE (host);


--
-- Name: hosts hosts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: guestbook guestbook_content_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guestbook
    ADD CONSTRAINT guestbook_content_type_fkey FOREIGN KEY (content_type) REFERENCES public.content_types(id);


--
-- Name: guestbook guestbook_host_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guestbook
    ADD CONSTRAINT guestbook_host_fkey FOREIGN KEY (host) REFERENCES public.hosts(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20241113175613');
