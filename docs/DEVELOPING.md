# developing

## configuring git

Run the following in the top level of the tree.

```
git config filter.sqlite3.clean 'sqlite3 %f .dump'
git config filter.sqlite3.smudge 'sqlite3 %f'
echo '*.db filter=sqlite3' >> .git/info/attributes
echo '*.sqlite filter=sqlite3' >> .git/info/attributes
```

You need to run this once on a clean checkout. Perhaps we can automate it somehow...