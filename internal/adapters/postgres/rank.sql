SELECT
    p.name,
    COUNT(*) as total_battles,
    SUM(CASE WHEN result = 'WIN' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN result = 'LOSE' THEN 1 ELSE 0 END) as losses,
    SUM(CASE WHEN result = 'ESCAPE' THEN 1 ELSE 0 END) as escape
FROM battle_logs bl
         JOIN players p on p.id = bl.player_id
WHERE created_at >= date_trunc('week', NOW())
group by p.name

---------------

SELECT
    p.name,
    m.name,
    count(*)
FROM battle_logs bl
         JOIN players p on p.id = bl.player_id
         JOIN monsters m on bl.monster_id = m.monster_id
group by m.name, p.name
order by p.name
