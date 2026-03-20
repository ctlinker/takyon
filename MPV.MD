# Minimal viable product

## Basic non-negociable

- Basic setup

```sh
takyon create black-arch -s 10000 -f ex4
takyon mount black-arch
takyon flash black-arch -d arch
takyon useradd black-arch -u mister --sudo
takyon enter black-arch -u mister
[mister@black-arch]#
```

- Regret option

```sh
takyon resize black-arch 15000
takyon remove back-arch
takyon fork black-arch white-arch
```

## More

```sh
takyon flash black-arch --env kde
takyon nested-server start # start an idle web server
takyon start black-arch --env kde --id 1 --dim 1280x720
```
