create table packages (
    url text primary key
);

create table package_versions (
    uuid uuid primary key default gen_random_uuid(),
    owner text not null references packages(url) on delete cascade,
    time text not null
);

create table internal_log (
    id text primary key check (id in ('last_fetch_time')),
    value text not null
);
