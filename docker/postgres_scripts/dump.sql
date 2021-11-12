create table if not exists profile(
  id serial not null primary key,
  create_time date default now(),
  update_time date default now(),
  name varchar(255) default '',
  email varchar(255) unique default '',
  password varchar(255) default '',
  date varchar(255) default '',
  description varchar(1000) default '',
  imgs varchar(255) [] default array [] :: varchar []
);

create trigger modify_payment_updated_at
    before update
    on profile
    for each row
execute procedure public.moddatetime(updated_at);

create table if not exists tag(
  id serial not null primary key,
  tag_name varchar(255) default ''
);
create table if not exists profile_tag(
  id serial not null primary key,
  profile_id integer,
  tag_id integer,
  constraint fk_pt_profile foreign key (profile_id) REFERENCES profile (id),
  constraint fk_pt_tag foreign key (tag_id) REFERENCES tag (id)
);
create table if not exists reactions(
  id serial not null primary key,
  id1 integer,
  id2 integer,
  type integer,
  constraint fk_pt_profile1 foreign key (id1) REFERENCES profile (id),
  constraint fk_pt_profile2 foreign key (id2) REFERENCES profile (id)
);
comment on column reactions.id1 is 'who like';
comment on column reactions.id2 is 'who liked';
create table if not exists matches(
  id serial not null primary key,
  id1 integer,
  id2 integer,
  constraint fk_pt_profile1 foreign key (id1) REFERENCES profile (id),
  constraint fk_pt_profile2 foreign key (id2) REFERENCES profile (id)
);
insert into
  tag(tag_name)
values('anime'),('music'),('gaming'),('sport'),('science');
