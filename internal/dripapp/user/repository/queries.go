package repository

const (
	GetUserQuery = `select id, email, password, name, gender, prefer, fromage, toage, date, 
	case when date <> '' then date_part('year', age(date::date)) else 0 end as age,
	description, imgs
	from profile where email = $1;`

	GetUserByIdAQuery = `select id, email, password, name, gender, prefer, fromage, toage, date, 
	case when date <> '' then date_part('year', age(date::date)) else 0 end as age,
	description, imgs
	from profile where id = $1;`

	CreateUserQuery = "INSERT into profile(email,password) VALUES($1,$2) RETURNING id, email, password;"

	UpdateUserQuery = `update profile set name=$2, gender=$3, prefer=$4, fromage=$5, toage=$6, date=$7, description=$8, imgs=$9 where email=$1
RETURNING id, email, password, name, gender, prefer, fromage, toage, date, 
case when date <> '' then date_part('year', age(date::date)) else 0 end as age, description, imgs;`

	DeleteTagsQuery = "delete from profile_tag where profile_id=$1 returning id;"

	GetTagsQuery = "select tagname from tag;"

	GetTagsByIdQuery = `select
							tagname
						from
							profile p
							join profile_tag pt on(pt.profile_id = p.id)
							join tag t on(pt.tag_id = t.id)
						where
							p.id = $1;`

	GetImgsByIDQuery = "SELECT imgs FROM profile WHERE id=$1;"

	InsertTagsQueryFirstPart = "insert into profile_tag(profile_id, tag_id) values"
	InsertTagsQueryParts     = "($1, (select id from tag where tagname=$%d))"

	UpdateImgsQuery = "update profile set imgs=$2 where id=$1 returning id;"

	AddReactionQuery = "insert into reactions(id1, id2, type) values ($1,$2,$3) returning id;"

	GetNextUserForSwipeQuery1 = `select op.id, op.email, op.password, op.name, op.gender, op.date,
									date_part('year', age(date::timestamp)) as age,
									op.description, op.reportstatus
									from profile op
									where op.id not in (
									select r.id2
									from reactions r
									where r.id1 = $1
									) and op.id not in (
									select m.id2
									from matches m
									where m.id1 = $1
									) and op.id not in (
									select r.id1
									from reactions r
									where r.id2 = $1 and type=0
									)
									and op.id <> $1
									and op.name <> ''
									and op.date <> ''
									and date_part('year', age(date::timestamp)) >= $2
									and date_part('year', age(date::timestamp)) <= $3;
									`

	GetNextUserForSwipeQueryPrefer = "and op.gender=$4\n"

	Limit = " limit 5;"

	GetUsersForMatchesQuery = `select
									op.id,
									op.email,
									op.password,
									op.name,
									op.date,
									case when op.date <> '' then date_part('year', age(op.date::timestamp)) else 0 end as age,
									op.description,
                  					op.reportstatus
								from profile p
								join matches m on (p.id = m.id1)
								join matches om on (om.id1 = m.id2 and om.id2 = m.id1)
								join profile op on (op.id = om.id1)
								where p.id = $1;`

	GetUsersForMatchesWithSearchingQuery = `select
												op.id,
												op.name,
												op.email,
												op.date,
												case when op.date <> '' then date_part('year', age(op.date::timestamp)) else 0 end as age,
												op.description,
                        						op.reportstatus
											from profile p
											join matches m on (p.id = m.id1)
											join matches om on (om.id1 = m.id2 and om.id2 = m.id1)
											join profile op on (op.id = om.id1)
											where p.id = $1 and LOWER(op.name) like LOWER($2);`

	GetLikesQuery = "select r.id1 from reactions r where r.id2 = $1 and r.type = 1;"

	DeleteReactionQuery = "delete from reactions r where ((r.id1=$1 and r.id2=$2) or (r.id1=$2 and r.id2=$1)) returning id;"

	DeleteMatchQuery = "delete from matches r where ((r.id1=$1 and r.id2=$2) or (r.id1=$2 and r.id2=$1)) returning id;"

	AddMatchQuery = "insert into matches(id1, id2) values ($1,$2),($2,$1) returning id;"

	GetUserLikesQuery = `select p.id,
							p.email,
							p.name,
							p.date,
							case when p.date <> '' then date_part('year', age(p.date::timestamp)) else 0 end as age,
							p.description,
							p.reportstatus
						from profile p
						join reactions r on
						(r.id1 = p.id
						and r.id2 = $1
						and r.type = 1
						and p.name <> ''
						and p.date <> '')
						where p.id not in (
						select r.id2
						from reactions r
						where r.id1 = $1 and type=0);
						`

	GetReportsQuery = "select reportdesc from reports;"

	GetReportIdFromDescQuery      = "select r.id from reports r where r.reportdesc = $1;"
	GetReportDescFromIdQuery      = "select reportdesc from reports r where r.id = $1;"
	AddReportToProfileQuery       = "insert into profile_report(profile_id, report_id) values($1, $2) returning id;"
	GetReportsCountQuery          = "select count(*) from profile_report where profile_id = $1;"
	GetReportsIdWithMaxCountQuery = `select report_id
									from profile_report
									group by report_id
									having count(*) = (
													select max(counts.c) from (
																				select count(*) c, report_id
																				from profile_report
																				where profile_id = $1
																				group by report_id
													) as counts);`
	UpdateProfilesReportStatusQuery = "update profile set reportstatus = $2 where id = $1 returning id;"

	CreatePaymentQuery = "insert into payment(id, status, amount, profile_id) values($1, $2, $3, $4) returning id;"
	UpdatePaymentQuery = "update payment set status=$2 where id=$1 returning id;"

	CreateSubscriptionQuery = "insert into subscription(period_start, period_end, profile_id, payment_id) values($1, $2, $3, $4) returning id;"
	UpdateSubscriptionQuery = "update subscription set paid=$2 where payment_id=$1 returning id;"

	CheckSubscriptionQuery = `select
								case when now() < period_end and paid then true
									 else false
								end
							  from subscription
							  where profile_id=$1
							  order by period_end DESC
							  limit 1;`
)
