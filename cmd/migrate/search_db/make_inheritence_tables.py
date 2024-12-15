from datetime import datetime
from datetime import timezone
import click
import json
import os

partitioned_tlds = set()
today = datetime.now(timezone.utc).strftime("%Y%m%d%H%M%S")
output_filename = f"db/migrations/{today}_partitions.sql"


def get_fp():
    fp = open(output_filename, "a")
    return fp


def safe(c):
    return c in "abcdefghijklmnopqrstuvwxyz0123456789_"


def partition_tld(fqdn, tld, tld_nibbles):
    clean = "".join(filter(safe, tld.lower()))
    if tld_nibbles not in partitioned_tlds:
        partitioned_tlds.add(tld_nibbles)
        fp = get_fp()
        fp.write("----------------------------------------\n")
        fp.write(f"-- {clean}\n")
        fp.write(" ".join(["-- ", tld_nibbles, "\n"]))
        fp.write("----------------------------------------\n")

        fp.write(
            f"create table if not exists {clean} () inherits (searchable_content);\n"
        )  # d64_{tld_nibbles}
        fp.write("\n")
        fp.close()


partitioned_domains = set()


def partition_domain(fqdn, rdomain, tld_nibbles, domain_nibbles):
    key = tld_nibbles + domain_nibbles
    tld = list(reversed(fqdn.split(".")))[0]
    clean = "".join(filter(safe, rdomain.lower()))
    if key not in partitioned_domains:
        partitioned_domains.add(key)
        fp = get_fp()
        fp.write("----------------------------------------\n")
        fp.write(f"-- {clean}\n")
        fp.write(" ".join(["-- ", tld_nibbles, domain_nibbles, "\n"]))
        fp.write("----------------------------------------\n")

        fp.write(
            f"create table if not exists {clean} () inherits ({tld});\n"
        )  # d64_{tld_nibbles}_{domain_nibbles}
        fp.write("\n")
        fp.close()


partitioned_subdomains = set()


def partition_subdomain(
    fqdn, rfqdn, rdomain, tld_nibbles, domain_nibbles, subdomain_nibbles
):
    key = tld_nibbles + domain_nibbles + subdomain_nibbles
    clean = "".join(filter(safe, rfqdn.lower()))
    clean_rdomain = "".join(filter(safe, rdomain.lower()))
    if key not in partitioned_subdomains:
        partitioned_subdomains.add(key)
        subdomain_plus_one = int(subdomain_nibbles, 16) + 1
        fp = get_fp()
        fp.write("----------------------------------------\n")
        fp.write(f"-- {clean}\n")
        fp.write(
            " ".join(["-- ", tld_nibbles, domain_nibbles, subdomain_nibbles, "\n"])
        )
        fp.write("----------------------------------------\n")

        fp.write(
            f"create table if not exists {clean} () inherits ({clean_rdomain});\n"  # d64_{tld_nibbles}_{domain_nibbles}_{subdomain_nibbles}
        )
        fp.write("\n")
        fp.close()


def migrate_up(tlds, jd64, start_int, end_int):
    fp = get_fp()
    fp.write("-- migrate:up\n\n")
    fp.close()
    for tld in tlds:
        d64tofqdn = jd64[tld]["Domain64ToFQDN"]
        for d64, fqdn in d64tofqdn.items():
            d64_int = int(d64, 16)
            if d64_int >= start_int and d64_int <= end_int:
                tld_nibbles = d64[0:2]
                domain_nibbles = d64[2:8]
                subdomain_nibbles = d64[8:14]
                tld = list(reversed(fqdn.split(".")))[0]
                rdomain = "_".join(list(reversed(fqdn.split(".")))[0:2])
                rfqdn = "_".join(list(reversed(fqdn.split("."))))
                partition_tld(fqdn, tld, tld_nibbles)
                partition_domain(fqdn, rdomain, tld_nibbles, domain_nibbles)
                partition_subdomain(
                    fqdn, rfqdn, rdomain, tld_nibbles, domain_nibbles, subdomain_nibbles
                )


def migrate_down(tlds, jd64, start_int, end_int):
    fp = get_fp()
    fp.write("-- migrate:down\n\n")
    fp.close()
    for tld in tlds:
        d64tofqdn = jd64[tld]["Domain64ToFQDN"]
        for d64, fqdn in d64tofqdn.items():
            d64_int = int(d64, 16)
            if d64_int >= start_int and d64_int <= end_int:
                tld_nibbles = d64[0:2]
                domain_nibbles = d64[2:8]
                subdomain_nibbles = d64[8:14]
                tld = list(reversed(fqdn.split(".")))[0]
                rdomain = "_".join(list(reversed(fqdn.split(".")))[0:2])
                rfqdn = "_".join(list(reversed(fqdn.split("."))))

                # if tld_nibbles in partitioned_tlds:
                #     fp = get_fp()
                #     fp.write(f"drop table if exists {tld};\n")
                #     partitioned_tlds.discard(tld_nibbles)
                #     fp.close()

                key = tld_nibbles + domain_nibbles
                if key in partitioned_domains:
                    clean = "".join(filter(safe, rdomain.lower()))
                    fp = get_fp()
                    fp.write(f"drop table if exists {clean};\n")
                    partitioned_domains.discard(key)
                    fp.close()

                key = tld_nibbles + domain_nibbles + subdomain_nibbles
                if key in partitioned_subdomains:
                    clean = "".join(filter(safe, rfqdn.lower()))
                    fp = get_fp()
                    fp.write(f"drop table if exists {clean};\n")
                    partitioned_subdomains.discard(key)
                    fp.close()


def trigger_function(tlds, jd64, start_int, end_int):
    fp = get_fp()
    fp.write("create or replace function searchable_content_insert_trigger_fun()\n")
    fp.write("  returns trigger as $$\n")
    fp.write("  begin\n")
    first = True
    for tld in tlds:
        d64tofqdn = jd64[tld]["Domain64ToFQDN"]
        for d64, fqdn in d64tofqdn.items():
            d64_int = int(d64, 16)
            if d64_int >= start_int and d64_int <= end_int:
                tld_nibbles = d64[0:2]
                domain_nibbles = d64[2:8]
                subdomain_nibbles = d64[8:14]
                subdomain_plus_one = int(subdomain_nibbles, 16) + 1
                tld = list(reversed(fqdn.split(".")))[0]
                rdomain = "_".join(list(reversed(fqdn.split(".")))[0:2])
                clean_rdomain = "".join(filter(safe, rdomain.lower()))
                rfqdn = "_".join(list(reversed(fqdn.split("."))))
                clean_rfqdn = "".join(filter(safe, rfqdn.lower()))
                d64_plus_one = (
                    tld_nibbles + domain_nibbles + f"{subdomain_plus_one:06x}"
                )
                if first:
                    cond = "if"
                    first = False
                else:
                    cond = "elsif"

                fp.write(
                    f"    {cond} (new.domain64 >= x'{d64}'::bigint and new.domain64 < x'{d64_plus_one}00'::bigint)\n"
                )
                fp.write(f"      then insert into {clean_rfqdn} values (new.*);\n")
    fp.write(f"    else insert into {tld} values (new.*);\n")
    fp.write("  end if;\n")
    fp.write("  return null;\n")
    fp.write("  end;\n")
    fp.write("$$\n")
    fp.write("language plpgsql;\n\n")
    fp.write("create trigger insert_searchable_content_trigger\n")
    fp.write("  before insert on searchable_content\n")
    fp.write(
        "  for each row execute function searchable_content_insert_trigger_fun();\n\n"
    )
    fp.close()


@click.command()
@click.argument("path")
@click.argument("migration_name")
@click.option("--start", default="0x0000000000000000", help="Start of partition range")
@click.option("--end", default="0x7FFFFFFFFFFFFF00", help="End of partition range")
def main(path, migration_name, start, end):
    global output_filename
    start_int = int(start, 16)
    end_int = int(end, 16)

    output_filename = f"db/migrations/{today}_inheritence_{migration_name}.sql"

    fp = open(path, "r")
    jd64 = json.load(fp)
    tlds = jd64["TLDs"]
    try:
        os.remove(output_filename)
    except:
        pass

    migrate_up(tlds, jd64, start_int, end_int)
    trigger_function(tlds, jd64, start_int, end_int)
    migrate_down(tlds, jd64, start_int, end_int)


if __name__ in "__main__":
    main()
