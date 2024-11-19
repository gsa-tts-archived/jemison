import requests
import click


@click.command()
@click.argument("filename")
def maine(filename):
    results = open("redirects.csv", "w")
    results.write("url,resolved,status\n")
    with open(filename) as file:
        for line in file:
            line = line.strip()
            try:
                r = requests.get("https://" + line, timeout=10)
                if r.history == []:
                    s = f"{line},{r.url},{r.status_code}"
                    print(s)
                    results.write(s + "\n")
                    results.flush()
                else:
                    for resp in r.history:
                        s = f"{line},{resp.url},{resp.status_code}"
                        print(s)
                        results.write(s + "\n")
                        results.flush()
                    s = f"{line},{r.url},{r.status_code}"
                    print(s)
                    results.write(s + "\n")
                    results.flush()

            except:
                s = f"{line},{r.url},0"
                print(s)


if __name__ in "__main__":
    maine()
