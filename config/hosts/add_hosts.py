import click
import csv
import json

# This takes a sitescanner-shaped CSV, and merges it into the appropriate 
# libsonnet documents in `domain64`. 

# target_url,
# base_domain,
# top_level_domain,
# branch,agency,
# bureau,
# source_list_federal_domains,
# source_list_dap,
# source_list_pulse,
# source_list_omb_idea,
# source_list_eotw,
# source_list_usagov,
# source_list_gov_man,
# source_list_uscourts,
# source_list_oira,
# source_list_other,
# source_list_mil_1,
# source_list_mil_2,
# omb_idea_public

@click.command()
@click.argument('filename')
def main(filename):
  

if __name__ in "__main__":
  main()