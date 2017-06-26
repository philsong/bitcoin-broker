--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE account (
    id integer NOT NULL,
    exchange character varying(64) NOT NULL,
    available_cny numeric(65,2) NOT NULL,
    available_btc numeric(65,4) NOT NULL,
    frozen_cny numeric(65,2) NOT NULL,
    frozen_btc numeric(65,4) NOT NULL,
    pause_trade boolean NOT NULL,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE account OWNER TO root;

--
-- Name: account_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE account_id_seq OWNER TO root;

--
-- Name: account_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE account_id_seq OWNED BY account.id;


--
-- Name: amount_config; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE amount_config (
    id integer NOT NULL,
    max_cny numeric(65,4) NOT NULL,
    max_btc numeric(65,4) NOT NULL,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE amount_config OWNER TO root;

--
-- Name: amount_config_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE amount_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE amount_config_id_seq OWNER TO root;

--
-- Name: amount_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE amount_config_id_seq OWNED BY amount_config.id;


--
-- Name: depth; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE depth (
    id integer NOT NULL,
    exchange character varying(64) NOT NULL,
    orderbook text NOT NULL,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE depth OWNER TO root;

--
-- Name: depth_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE depth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE depth_id_seq OWNER TO root;

--
-- Name: depth_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE depth_id_seq OWNED BY depth.id;


--
-- Name: django_migrations; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE django_migrations (
    id integer NOT NULL,
    app character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    applied timestamp with time zone NOT NULL
);


ALTER TABLE django_migrations OWNER TO root;

--
-- Name: django_migrations_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE django_migrations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE django_migrations_id_seq OWNER TO root;

--
-- Name: django_migrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE django_migrations_id_seq OWNED BY django_migrations.id;


--
-- Name: exchange_config; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE exchange_config (
    id integer NOT NULL,
    exchange character varying(64) NOT NULL,
    access_key character varying(64) NOT NULL,
    secret_key character varying(64) NOT NULL
);


ALTER TABLE exchange_config OWNER TO root;

--
-- Name: exchange_config_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE exchange_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE exchange_config_id_seq OWNER TO root;

--
-- Name: exchange_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE exchange_config_id_seq OWNED BY exchange_config.id;


--
-- Name: site_order; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE site_order (
    created timestamp with time zone DEFAULT now() NOT NULL,
    id integer NOT NULL,
    client_id character varying(64) NOT NULL,
    trade_type character varying(32) NOT NULL,
    order_status character varying(32) NOT NULL,
    amount numeric(65,4) NOT NULL,
    estimate_price numeric(65,2) NOT NULL,
    estimate_cny numeric(65,2) NOT NULL,
    estimate_btc numeric(65,4) NOT NULL,
    price numeric(65,2) NOT NULL
);


ALTER TABLE site_order OWNER TO root;

--
-- Name: site_order_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE site_order_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE site_order_id_seq OWNER TO root;

--
-- Name: site_order_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE site_order_id_seq OWNED BY site_order.id;


--
-- Name: ticker; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE ticker (
    id integer NOT NULL,
    ask numeric(65,2) NOT NULL,
    bid numeric(65,2) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE ticker OWNER TO root;

--
-- Name: ticker_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE ticker_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE ticker_id_seq OWNER TO root;

--
-- Name: ticker_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE ticker_id_seq OWNED BY ticker.id;


--
-- Name: trade_order; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE trade_order (
    created timestamp with time zone DEFAULT now() NOT NULL,
    id integer NOT NULL,
    exchange character varying(64) NOT NULL,
    trade_type character varying(32) NOT NULL,
    order_status character varying(32) NOT NULL,
    estimate_cny numeric(65,2) NOT NULL,
    estimate_btc numeric(65,4) NOT NULL,
    estimate_price numeric(65,2) NOT NULL,
    deal_cny numeric(65,2) NOT NULL,
    deal_btc numeric(65,4) NOT NULL,
    deal_price numeric(65,2) NOT NULL,
    order_id character varying(64) NOT NULL,
    info text NOT NULL,
    site_order_id integer NOT NULL,
    price_margin numeric(65,4) NOT NULL,
    match_id integer NOT NULL,
    memo text NOT NULL,
    price numeric(65,2) NOT NULL,
    try_times integer NOT NULL,
    update_at timestamp with time zone NOT NULL
);


ALTER TABLE trade_order OWNER TO root;

--
-- Name: trade_order_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE trade_order_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE trade_order_id_seq OWNER TO root;

--
-- Name: trade_order_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE trade_order_id_seq OWNED BY trade_order.id;


--
-- Name: trader_traderbusevent; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE trader_traderbusevent (
    id integer NOT NULL,
    uuid character varying(32) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    visited_at timestamp with time zone NOT NULL,
    name character varying(200) NOT NULL,
    has_pushed boolean NOT NULL,
    content_type_id integer NOT NULL,
    object_id character varying(200) NOT NULL,
    data text NOT NULL,
    request_info text NOT NULL,
    CONSTRAINT trader_traderbusevent_content_type_id_check CHECK ((content_type_id >= 0))
);


ALTER TABLE trader_traderbusevent OWNER TO root;

--
-- Name: trader_traderbusevent_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE trader_traderbusevent_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE trader_traderbusevent_id_seq OWNER TO root;

--
-- Name: trader_traderbusevent_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE trader_traderbusevent_id_seq OWNED BY trader_traderbusevent.id;


--
-- Name: trader_traderbusposition; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE trader_traderbusposition (
    id integer NOT NULL,
    uuid character varying(32) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    visited_at timestamp with time zone NOT NULL,
    last_pos integer NOT NULL,
    CONSTRAINT trader_traderbusposition_last_pos_check CHECK ((last_pos >= 0))
);


ALTER TABLE trader_traderbusposition OWNER TO root;

--
-- Name: trader_traderbusposition_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE trader_traderbusposition_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE trader_traderbusposition_id_seq OWNER TO root;

--
-- Name: trader_traderbusposition_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE trader_traderbusposition_id_seq OWNED BY trader_traderbusposition.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY account ALTER COLUMN id SET DEFAULT nextval('account_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY amount_config ALTER COLUMN id SET DEFAULT nextval('amount_config_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY depth ALTER COLUMN id SET DEFAULT nextval('depth_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY django_migrations ALTER COLUMN id SET DEFAULT nextval('django_migrations_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY exchange_config ALTER COLUMN id SET DEFAULT nextval('exchange_config_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY site_order ALTER COLUMN id SET DEFAULT nextval('site_order_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY ticker ALTER COLUMN id SET DEFAULT nextval('ticker_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY trade_order ALTER COLUMN id SET DEFAULT nextval('trade_order_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY trader_traderbusevent ALTER COLUMN id SET DEFAULT nextval('trader_traderbusevent_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY trader_traderbusposition ALTER COLUMN id SET DEFAULT nextval('trader_traderbusposition_id_seq'::regclass);


--
-- Name: account_exchange_key; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_exchange_key UNIQUE (exchange);


--
-- Name: account_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: amount_config_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY amount_config
    ADD CONSTRAINT amount_config_pkey PRIMARY KEY (id);


--
-- Name: depth_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY depth
    ADD CONSTRAINT depth_pkey PRIMARY KEY (id);


--
-- Name: django_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY django_migrations
    ADD CONSTRAINT django_migrations_pkey PRIMARY KEY (id);


--
-- Name: exchange_config_exchange_key; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY exchange_config
    ADD CONSTRAINT exchange_config_exchange_key UNIQUE (exchange);


--
-- Name: exchange_config_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY exchange_config
    ADD CONSTRAINT exchange_config_pkey PRIMARY KEY (id);


--
-- Name: site_order_client_id_key; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY site_order
    ADD CONSTRAINT site_order_client_id_key UNIQUE (client_id);


--
-- Name: site_order_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY site_order
    ADD CONSTRAINT site_order_pkey PRIMARY KEY (id);


--
-- Name: ticker_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY ticker
    ADD CONSTRAINT ticker_pkey PRIMARY KEY (id);


--
-- Name: trade_order_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY trade_order
    ADD CONSTRAINT trade_order_pkey PRIMARY KEY (id);


--
-- Name: trader_traderbusevent_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY trader_traderbusevent
    ADD CONSTRAINT trader_traderbusevent_pkey PRIMARY KEY (id);


--
-- Name: trader_traderbusposition_pkey; Type: CONSTRAINT; Schema: public; Owner: root; Tablespace: 
--

ALTER TABLE ONLY trader_traderbusposition
    ADD CONSTRAINT trader_traderbusposition_pkey PRIMARY KEY (id);


--
-- Name: account_afd1a1a8; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX account_afd1a1a8 ON account USING btree (updated_at);


--
-- Name: account_exchange_cfd744d3_like; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX account_exchange_cfd744d3_like ON account USING btree (exchange varchar_pattern_ops);


--
-- Name: account_fde81f11; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX account_fde81f11 ON account USING btree (created_at);


--
-- Name: amount_config_afd1a1a8; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX amount_config_afd1a1a8 ON amount_config USING btree (updated_at);


--
-- Name: amount_config_fde81f11; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX amount_config_fde81f11 ON amount_config USING btree (created_at);


--
-- Name: depth_afd1a1a8; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX depth_afd1a1a8 ON depth USING btree (updated_at);


--
-- Name: depth_fde81f11; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX depth_fde81f11 ON depth USING btree (created_at);


--
-- Name: exchange_config_exchange_ff3dc00a_like; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX exchange_config_exchange_ff3dc00a_like ON exchange_config USING btree (exchange varchar_pattern_ops);


--
-- Name: site_order_client_id_80f1b664_like; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX site_order_client_id_80f1b664_like ON site_order USING btree (client_id varchar_pattern_ops);


--
-- Name: ticker_created_at_a5a68202_uniq; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX ticker_created_at_a5a68202_uniq ON ticker USING btree (created_at);


--
-- Name: ticker_updated_at_e72cc59d_uniq; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX ticker_updated_at_e72cc59d_uniq ON ticker USING btree (updated_at);


--
-- Name: trade_order_2637814c; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trade_order_2637814c ON trade_order USING btree (site_order_id);


--
-- Name: trader_traderbusevent_4c53479a; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trader_traderbusevent_4c53479a ON trader_traderbusevent USING btree (has_pushed);


--
-- Name: trader_traderbusevent_afd1a1a8; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trader_traderbusevent_afd1a1a8 ON trader_traderbusevent USING btree (updated_at);


--
-- Name: trader_traderbusevent_fde81f11; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trader_traderbusevent_fde81f11 ON trader_traderbusevent USING btree (created_at);


--
-- Name: trader_traderbusposition_afd1a1a8; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trader_traderbusposition_afd1a1a8 ON trader_traderbusposition USING btree (updated_at);


--
-- Name: trader_traderbusposition_fde81f11; Type: INDEX; Schema: public; Owner: root; Tablespace: 
--

CREATE INDEX trader_traderbusposition_fde81f11 ON trader_traderbusposition USING btree (created_at);


--
-- Name: trade_order_site_order_id_6f9e5a14_fk_site_order_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY trade_order
    ADD CONSTRAINT trade_order_site_order_id_6f9e5a14_fk_site_order_id FOREIGN KEY (site_order_id) REFERENCES site_order(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: public; Type: ACL; Schema: -; Owner: root
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM root;
GRANT ALL ON SCHEMA public TO root;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

