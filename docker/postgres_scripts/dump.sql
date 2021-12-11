CREATE extension IF NOT EXISTS citext;
create table if not exists profile(
  id serial not null primary key,
  create_time timestamp default now(),
  update_time timestamp default now(),
  email citext,
  password varchar(100) default '',
  name varchar(63) default '',
  gender varchar(15) default '',
  prefer varchar(15) default '',
  fromage smallint default 18,
  toage smallint default 100,
  date varchar(15) default '',
  description varchar(1023) default '',
  imgs varchar(255) [] default array [] :: varchar [],
  reportstatus varchar(255) default ''
);
create table if not exists tag(
  id serial not null primary key,
  tagname varchar(255) default ''
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
create table message(
  message_id serial not null primary key,
  from_id integer,
  to_id integer,
  text text default '',
  date timestamptz default now(),
  constraint fk_ms_profile1 foreign key (from_id) REFERENCES profile (id),
  constraint fk_ms_profile2 foreign key (to_id) REFERENCES profile (id)
);
-- message date index
  create index idx_ms_date on message(date) include (from_id, to_id, text);

create table if not exists reports(
  id serial not null primary key,
  reportdesc varchar(255) default ''
);
create table if not exists profile_report(
  id serial not null primary key,
  profile_id integer,
  report_id integer,
  constraint fk_pr_profile foreign key (profile_id) REFERENCES profile (id),
  constraint fk_pr_report foreign key (report_id) REFERENCES reports (id)
);
create table if not exists payment(
  id serial not null primary key,
  period timestamptz default now(),
  status smallint,
  profile_id integer,
  constraint fk_pm_profile foreign key (profile_id) REFERENCES profile (id)
);


insert into
  tag(tagname)
values('аниме'),('рок'),('игры'),('спорт'),('наука'),('рэп'),('джаз'),('западная музыка'),('комедии'),('футбол');
insert into
  reports(reportdesc)
values('Фалишивый профиль/спам'),('Непристойное общение'),('Скам'),('Несовершеннолетний пользователь');

-- foregn keys
  create index idx_ms_from_id on message(from_id);
create index idx_ms_to_id on message(to_id);
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
returns trigger
as $$
  begin
    NEW.update_time = NOW();
return NEW;
end;
$$ language plpgsql;
create trigger modify_payment_update_time before
update
  on profile for each row execute procedure moddatetime();
vacuum full analyze;