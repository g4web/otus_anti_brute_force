CREATE TYPE network_type AS ENUM ('white', 'black');

create table network
(
    id      serial
        constraint network_pk
            primary key,
    network varchar(18),
    type    network_type not null
);

comment on column network.network is '192.168.0.0/24'
