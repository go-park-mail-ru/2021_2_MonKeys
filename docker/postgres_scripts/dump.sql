CREATE extension IF NOT EXISTS citext;

create table if not exists profile(
  id serial not null primary key,
  create_time timestamp default now(),
  update_time timestamp default now(),
  name varchar(255) default '',
  email citext,
  password varchar(255) default '',
  date varchar(255) default '',
  description varchar(1000) default '',
  imgs varchar(255) [] default array [] :: varchar []
);

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

-- foregn keys
create index idx_react_id1 on reactions(id1);
create index idx_react_id2 on reactions(id2);
create index idx_match_id1 on matches(id1);
create index idx_match_id2 on matches(id1);
create index idx_pt_profile_id on profile_tag(profile_id);
create index idx_pt_tag_id on profile_tag(tag_id);

-- lower
create unique index uniq_email ON profile(email);

-- search
create index idx_profile_imgs_gin on profile using gin (imgs);

create or replace function moddatetime()
returns trigger as $$
begin
  NEW.update_time = NOW();
  return NEW;
end;
$$ language plpgsql;

create trigger modify_payment_update_time
    before update
    on profile
    for each row
execute procedure moddatetime();

vacuum full analyze;