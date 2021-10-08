create table profile(  
    id serial not null primary key,
    create_time date default now(),
    update_time date default  now(),
    name varchar(255) null,
    email varchar(255) null,
    password varchar(255) null,
    date date null,
    description varchar(255) null,
    img varchar(255) null
);

create table tag(  
    id serial not null primary key,
    tag_name varchar(255) null
);

create table profile_tag(  
    id serial not null primary key,
    profile_id integer,
    tag_id integer, 
    constraint fk_pt_profile
     foreign key (profile_id) 
     REFERENCES profile (id),
    constraint fk_pt_tag
     foreign key (tag_id) 
     REFERENCES tag (id)
);

create table reactions(  
    id serial not null primary key,
    id1 integer,
    id2 integer,
    type varchar(255) null,
    constraint fk_pt_profile1
     foreign key (id1) 
     REFERENCES profile (id),
    constraint fk_pt_profile2
     foreign key (id2) 
     REFERENCES profile (id)
);
comment on column likes.id1 is 'who like';
comment on column likes.id2 is 'who liked';

create table matches(  
    id serial not null primary key,
    id1 integer,
    id2 integer, 
    constraint fk_pt_profile1
     foreign key (id1) 
     REFERENCES profile (id),
    constraint fk_pt_profile2
     foreign key (id2) 
     REFERENCES profile (id)
);