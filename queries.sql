-- /group/:id
-- Get next week's prompt
select wp.name, wp.week, wp.spotify_url, u.name
from weekly_prompts wp
join group_members gm on wp.creator_id = gm.user_id 
join users u on wp.creator_id = u.id
where gm.group_id = $1
and wp.deadline > now()
order by wp.datetime asc
limit 1;

-- /group/:id
-- Get all done prompts
select wp.id, wp.name, wp.spotify_url, u.name
from weekly_prompts wp
join group_members gm on wp.creator_id = gm.user_id
join users u on wp.creator_id = u.id
where gm.group_id = $1
and wp.deadline < now()
order by wp.datetime desc;

-- /group/:id
-- Get upcoming prompts
select wp.id, wp.name, wp.spotify_url
from weekly_prompts wp
join group_members gm on wp.creator_id = gm.user_id
where gm.group_id = $1
and wp.deadline > now()
order by wp.datetime asc;


-- /group/:id/week/:num
-- Get prompt details
select wp.id, wp.name, wp.deadline, t.spotify_url, t.title, t.artist, t.album_cover
from weekly_prompts wp
join group_members gm on wp.creator_id = gm.user_id
join picks p on p.weekly_prompt_id = wp.id
where gm.group_id = $1
and wp.week = $2;

-- 
