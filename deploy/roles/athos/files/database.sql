CREATE database dartagnan;
CREATE USER dtuser WITH PASSWORD 'SECRET';
ALTER USER dtuser WITH SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE dartagnan to dtuser;
\connect dartagnan

CREATE TABLE subscription_plans (
    id serial not null primary key,
    code character varying(1024) not null,
    name character varying(1024) not null,
    description character varying(1024) not null,
    price numeric,
    period numeric default null
);

INSERT INTO subscription_plans VALUES (1, 'trial', 'Trial Pizza', '30 Day Trial', 0, 30);
INSERT INTO subscription_plans VALUES (2, 'crostino', 'Crostino', '- Stable Updates repository 
- Community Support
- Support tickets not included / 100 € each', 48.00, 365);
INSERT INTO subscription_plans VALUES (3, 'lasagna', 'Lasagna', '- Stable Updates repository 
- Professional support via Email + SSH
- 3 support tickets/year included', 250.0, 365);
INSERT INTO subscription_plans VALUES (4, 'fiorentina', 'Fiorentina', '- Stable Updates repository 
- Professional support via Email + SSH
- 6 support tickets/year included
- Monitoring Portal', 450.0, 365);
INSERT INTO subscription_plans VALUES (5, 'pizza', 'Pizza', '- Stable Updates repository 
- Professional support via Email + SSH + Phone
- 12 support tickets/year included
- Monitoring Portal', 800.0, 365);

CREATE TABLE subscriptions (
    id serial not null primary key,
    user_id character varying(1024) not null,
    subscription_plan_id bigint not null references subscription_plans(id),
    valid_from timestamp default current_timestamp,
    valid_until timestamp default current_timestamp,
    status character varying(10) default null,
    created timestamp default current_timestamp
);

CREATE TABLE systems (
    id serial not null primary key,
    subscription_id bigint not null references subscriptions(id),
    creator_id character varying(1024) not null,
    uuid character varying(1024) not null,
    secret character varying(1024) not null,
    tags character varying(1024) null,
    public_ip character varying(1024) null,
    status character varying(1024) default null,
    created timestamp default current_timestamp,
    notification jsonb,
    UNIQUE(uuid),
    UNIQUE(secret)
);

CREATE TABLE inventories (
  id serial not null primary key,
  system_id bigint not null references systems(id) on delete cascade,
  timestamp timestamp default current_timestamp,
  data jsonb,
  UNIQUE(system_id)
);

CREATE TABLE inventory_histories (
  id serial not null primary key,
  system_id bigint not null references systems(id),
  timestamp timestamp default current_timestamp,
  data jsonb
);

CREATE TABLE alerts (
  id serial not null primary key,
  system_id bigint not null references systems(id) on delete cascade,
  alert_id character varying(1024) not null,
  timestamp timestamp default current_timestamp,
  status character varying(1024) default null,
  note character varying(1024) not null,
  priority character varying(1024) default 0
);

CREATE TABLE alert_histories (
  id serial not null primary key,
  system_id bigint not null references systems(id) on delete cascade,
  alert_id character varying(1024) not null,
  start_time timestamp default null,
  end_time timestamp default null,
  status_from character varying(1024) default null,
  status_to character varying(1024) default null,
  priority character varying(1024) default 0,
  resolution character varying(1024) default 0
);

CREATE TABLE heartbeats (
  id serial not null primary key,
  system_id bigint not null references systems(id) on delete cascade,
  timestamp timestamp default current_timestamp,
  UNIQUE(system_id)
);

CREATE TABLE payments (
  id serial not null primary key,
  creator_id character varying(1024) not null,
  payment character varying(1024) not null,
  system_id bigint not null,
  created timestamp default current_timestamp,
  UNIQUE(payment)
);

CREATE TABLE taxes (
  country character varying(1024) not null primary key,
  percentage integer
);

CREATE TABLE billings (
  id serial not null primary key,
  creator_id character varying(1024) not null,
  name character varying(1024) not null,
  address character varying(1024) not null,
  city character varying(1024) not null,
  postal_code character varying(1024) not null,
  country character varying(1024) not null references taxes(country),
  vat character varying(1024),
  UNIQUE(creator_id)
);

INSERT INTO taxes VALUES ('Other',0);
INSERT INTO taxes VALUES ('Belgium',21);
INSERT INTO taxes VALUES ('Bulgaria',20);
INSERT INTO taxes VALUES ('Czech Republic',21);
INSERT INTO taxes VALUES ('Denmark',25);
INSERT INTO taxes VALUES ('Germany',19);
INSERT INTO taxes VALUES ('Estonia',20);
INSERT INTO taxes VALUES ('Ireland',23);
INSERT INTO taxes VALUES ('Greece',24);
INSERT INTO taxes VALUES ('Spain',21);
INSERT INTO taxes VALUES ('France',20);
INSERT INTO taxes VALUES ('Croatia',25);
INSERT INTO taxes VALUES ('Italy',22);
INSERT INTO taxes VALUES ('Cyprus',19);
INSERT INTO taxes VALUES ('Latvia',21);
INSERT INTO taxes VALUES ('Lithuania',21);
INSERT INTO taxes VALUES ('Luxembourg',17);
INSERT INTO taxes VALUES ('Hungary',27);
INSERT INTO taxes VALUES ('Malta',18);
INSERT INTO taxes VALUES ('Netherlands',21);
INSERT INTO taxes VALUES ('Austria',20);
INSERT INTO taxes VALUES ('Poland',23);
INSERT INTO taxes VALUES ('Portugal',23);
INSERT INTO taxes VALUES ('Romania',19);
INSERT INTO taxes VALUES ('Slovenia',22);
INSERT INTO taxes VALUES ('Slovakia',20);
INSERT INTO taxes VALUES ('Finland',24);
INSERT INTO taxes VALUES ('Sweden',25);
INSERT INTO taxes VALUES ('United Kingdom',20);
