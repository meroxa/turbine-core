--
-- PostgreSQL database dump
--

ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_pkey;
ALTER TABLE public.orders ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE public.orders_id_seq;
DROP TABLE public.orders;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    category character varying(255) NOT NULL,
    product_type character varying(255) NOT NULL,
    product_name text NOT NULL,
    stock boolean DEFAULT true,
    product_id bigint NOT NULL,
    shipping_address text NOT NULL,
    customer_email text NOT NULL
);

ALTER TABLE ONLY public.orders REPLICA IDENTITY FULL;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);

--
-- PostgreSQL database dump complete
--