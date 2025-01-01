import bs4
import click
import requests

@click.command()
@click.argument('url')
@click.argument('page_name')
def main(url, page_name):
  r = requests.get(url)
  url_contents = r.text
  soup = bs4.BeautifulSoup(url_contents, features="lxml")
  contents = soup.getText()
  with open(f"static/{page_name}.html", 'w') as fp:
    fp.write("<html>\n")
    fp.write(f"<head><title>{page_name}</title></head>\n")
    fp.write("<body>\n")
    fp.write(contents.strip())
    fp.write("\n</body></html>\n")

if __name__ in "__main__":
  main()