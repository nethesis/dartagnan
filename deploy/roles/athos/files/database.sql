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
INSERT INTO subscription_plans VALUES (1, 'trial', 'Trial', 'Trial of 30d', 0.00, 30);
INSERT INTO subscription_plans VALUES (2, 'crostino', 'Crostino', 'Good Starter', 48.00, 360);
INSERT INTO subscription_plans VALUES (3, 'lasagna', 'Lasagna', 'Homemade first plate', 250.00, 360);
INSERT INTO subscription_plans VALUES (4, 'fiorentina', 'Fiorentina', 'The main course', 450.00, 360);
INSERT INTO subscription_plans VALUES (5, 'pizza', 'Pizza', 'What’s else?', 800.00, 360);

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

CREATE TABLE billings (
  id serial not null primary key,
  creator_id character varying(1024) not null,
  name character varying(1024) not null,
  address character varying(1024) not null,
  nation character varying(1024) not null,
  type character varying(1024) not null,
  vat character varying(1024),
  UNIQUE(creator_id)
);

