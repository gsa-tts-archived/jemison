# generating partition tables

For all of the .gov domains:

python make_inheritance_tables.py \
  ../../../config/domain64/domain64.json \
  all_dot_gov \
  --start 0x0100000000000000 \
  --end 0x01FF000000000000


For NIH:

python make_inheritance_tables.py \
  ../../../config/domain64/domain64.json \
  just_nih \
  --start 0x0100008D00000100 \
  --end 0x0100008DFFFFFF00

