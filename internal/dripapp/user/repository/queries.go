package repository

const (
	GetUserQuery = "select id, name, email, password, date, description, imgs from profile where email = $1;"

	GetUserByIdAQuery = "select id, name, email, password, date, description, imgs from profile where id = $1;"

	CreateUserQuery = "INSERT into profile(email,password) VALUES($1,$2) RETURNING id, email, password;"

	UpdateUserQuery = "update profile set name=$1, date=$3, description=$4, imgs=$5 where email=$2 RETURNING id, name, email, password, date, description, imgs;"

	DeleteTagsQuery = "delete from profile_tag where profile_id=$1 returning id;"

	GetTagsQuery = "select tag_name from tag;"

	GetTagsByIdQuery = `select
							tag_name
						from
							profile p
							join profile_tag pt on(pt.profile_id = p.id)
							join tag t on(pt.tag_id = t.id)
						where
							p.id = $1;`

	GetImgsByIDQuery = "SELECT imgs FROM profile WHERE id=$1;"

	InsertTagsQueryFirstPart = "insert into profile_tag(profile_id, tag_id) values"
	InsertTagsQueryParts     = "($1, (select id from tag where tag_name=$%d))"

	UpdateImgsQuery = "update profile set imgs=$2 where id=$1 returning id;"

	AddReactionQuery = "insert into reactions(id1, id2, type) values ($1,$2,$3) returning id;"

	GetNextUserForSwipeQuery = `select 
									op.id,
									op.name,
									op.email,
									op.password,
									op.date,
									op.description
								from profile op
								where op.id not in (
									select r.id2
									from reactions r
									where r.id1 = $1
								) and op.id not in (
									select m.id2
									from matches m
									where m.id1 = $1
								) and op.id <> $1
								and op.name <> ''
								and op.age <> ''
									limit 5;`

	GetUsersForMatchesQuery = `select
									op.id,
									op.name,
									op.email,
									op.password,
									op.date,
									op.description
								from profile p
								join matches m on (p.id = m.id1)
								join matches om on (om.id1 = m.id2 and om.id2 = m.id1)
								join profile op on (op.id = om.id1)
								where p.id = $1;`

	GetLikesQuery = "select r.id1 from reactions r where r.id2 = $1 and r.type = 1;"

	DeleteLikeQuery = "delete from reactions r where ((r.id1=$1 and r.id2=$2) or (r.id1=$2 and r.id2=$1)) returning id;"

	AddMatchQuery = "insert into matches(id1, id2) values ($1,$2),($2,$1) returning id;"
)
