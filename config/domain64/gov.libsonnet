local assertion = import 'assertions.libsonnet';
local util = import 'util.libsonnet';

local domains = {
  '000001': {
    children: {
      '000001': 'agile',
      '000002': 'derisking-guide',
      '000003': 'accessibility',
      '000004': 'before-you-ship',
      '000005': 'ux-guide',
      '000006': 'product-guide',
      '000007': 'engineering',
      '000008': 'content-guide',
      '000009': 'methods',
      '00000A': 'guides',
    },
    name: '18f',
  },
  '000002': {
    children: {
      '000001': 'www',
    },
    name: '9-11commission',
  },
  '000003': {
    children: {
      '000001': 'ictbaseline',
      '000002': 'beta',
      '000003': 'www',
    },
    name: 'access-board',
  },
  '000004': {
    children: {
      '000001': 'repatriation',
    },
    name: 'acf',
  },
  '000005': {
    children: {
      '000001': 'beta',
      '000002': 'archive',
      '000003': 'www',
    },
    name: 'ada',
  },
  '000006': {
    children: {},
    name: 'adlnet',
  },
  '000007': {
    children: {
      '000001': 'publications',
      '000002': 'share',
    },
    name: 'america',
  },
  '000008': {
    children: {},
    name: 'americorps',
  },
  '000009': {
    children: {
      '000001': 'www',
    },
    name: 'apprenticeship',
  },
  '00000A': {
    children: {
      '000001': 'www',
      '000002': 'obamawhitehouse',
      '000003': 'founders',
      '000004': 'georgewbush-whitehouse',
      '000005': 'situationroom',
      '000006': 'obamawhitehouse.open',
      '000007': 'blogs.reagan',
      '000008': 'blogs.isoo',
      '000009': 'blogs.declassification',
      '00000A': 'blogs.transforming-classification',
      '00000B': 'hoover',
      '00000C': 'blogs.annotation',
      '00000D': 'blogs.records-express',
      '00000E': 'blogs.foia',
      '00000F': 'blogs.hoover',
      '000010': 'museum',
      '000011': 'blogs.jfk',
      '000012': 'blogs.rediscovering-black-history',
      '000013': 'blogs.education',
      '000014': 'blogs.narations',
      '000015': 'blogs.aotus',
      '000016': 'blogs.fdr',
      '000017': 'obamawhitehouse.letsmove',
      '000018': 'blogs.unwritten-record',
      '000019': 'blogs.text-message',
      '00001A': 'catalog',
      '00001B': 'blogs.prologue',
      '00001C': 'trumpwhitehouse',
      '00001D': 'clintonwhitehouse3',
      '00001E': 'clintonwhitehouse4',
      '00001F': 'clintonwhitehouse6',
    },
    name: 'archives',
  },
  '00000B': {
    children: {
      '000001': 'www',
    },
    name: 'atf',
  },
  '00000C': {
    children: {
      '000001': 'ssabest',
      '000002': 'www',
    },
    name: 'benefits',
  },
  '00000D': {
    children: {
      '000001': 'www',
    },
    name: 'bep',
  },
  '00000E': {
    children: {
      '000002': 'www',
    },
    name: 'bjs',
  },
  '00000F': {
    children: {
      '000001': 'www',
    },
    name: 'blm',
  },
  '000010': {
    children: {
      '000001': 'www',
    },
    name: 'boem',
  },
  '000011': {
    children: {
      '000001': 'www',
    },
    name: 'bop',
  },
  '000012': {
    children: {
      '000001': 'ntl.rosap',
    },
    name: 'bts',
  },
  '000013': {
    children: {},
    name: 'buildbackbetter',
  },
  '000014': {
    children: {
      '000001': 'datascience',
      '000002': 'cancercontrol.ebccp',
      '000003': 'www',
    },
    name: 'cancer',
  },
  '000015': {
    children: {},
    name: 'cbca',
  },
  '000016': {
    children: {
      '000001': 'biometrics.www',
      '000002': 'www',
    },
    name: 'cbp',
  },
  '000017': {
    children: {},
    name: 'ccb',
  },
  '000018': {
    children: {
      '000001': 'www',
    },
    name: 'cdc',
  },
  '000019': {
    children: {
      '000001': 'www',
    },
    name: 'cdfifund',
  },
  '00001A': {
    children: {
      '000001': 'www',
    },
    name: 'cdo',
  },
  '00001B': {
    children: {
      '000001': 'www',
    },
    name: 'census',
  },
  '00001C': {
    children: {
      '000001': 'www',
    },
    name: 'cfa',
  },
  '00001D': {
    children: {
      '000001': 'www',
    },
    name: 'cfo',
  },
  '00001E': {
    children: {
      '000001': 'www',
    },
    name: 'challenge',
  },
  '00001F': {
    children: {
      '000001': 'www',
    },
    name: 'cia',
  },
  '000020': {
    children: {
      '000001': 'tmf',
      '000002': 'www',
    },
    name: 'cio',
  },
  '000021': {
    children: {
      '000001': 'us-cert',
      '000002': 'www',
      '000003': 'niccs',
    },
    name: 'cisa',
  },
  '000022': {
    children: {
      '000001': 'www',
    },
    name: 'citizenscience',
  },
  '000023': {
    children: {
      '000001': 'toolkit',
      '000002': 'www',
    },
    name: 'climate',
  },
  '000024': {
    children: {
      '000000': '',
      '000001': 'app.fec-prod-proxy',
    },
    name: 'cloud',
  },
  '000025': {
    children: {
      '000001': 'www',
      '000002': 'partnershipforpatients',
      '000003': 'qpp',
      '000004': 'regulations-pilot',
      '000005': 'innovation',
    },
    name: 'cms',
  },
  '000026': {
    children: {
      '000001': 'www',
    },
    name: 'cmts',
  },
  '000027': {
    children: {},
    name: 'code',
  },
  '000028': {
    children: {
      '000001': 'www',
    },
    name: 'coldcaserecords',
  },
  '000029': {
    children: {
      '000001': 'www',
    },
    name: 'collegedrinkingprevention',
  },
  '00002A': {
    children: {
      '000001': '2017-2021',
      '000002': '2014-2017',
      '000003': '2010-2014',
      '000004': 'www',
    },
    name: 'commerce',
  },
  '00002B': {
    children: {
      '000001': 'beta',
      '000002': 'www',
    },
    name: 'consumerfinance',
  },
  '00002C': {
    children: {
      '000001': 'www',
    },
    name: 'copyright',
  },
  '00002D': {
    children: {},
    name: 'coralreef',
  },
  '00002E': {
    children: {
      '000001': 'www',
    },
    name: 'crimesolutions',
  },
  '00002F': {
    children: {
      '000002': 'www',
    },
    name: 'cttso',
  },
  '000030': {
    children: {
      '000001': 'www',
    },
    name: 'cuidadodesalud',
  },
  '000031': {
    children: {
      '000001': 'resources',
    },
    name: 'data',
  },
  '000032': {
    children: {
      '000001': 'www',
    },
    name: 'dataprivacyframework',
  },
  '000033': {
    children: {
      '000001': 'cfsadashboard',
      '000002': 'rhc',
      '000003': 'dcra',
      '000004': 'dhcf',
      '000005': 'abra',
    },
    name: 'dc',
  },
  '000034': {
    children: {
      '000001': 'www',
      '000002': 'basicresearch',
      '000003': 'media',
      '000004': 'minerva',
      '000005': 'prhome',
      '000006': 'dod',
      '000007': 'dpcld',
      '000008': 'comptroller',
    },
    name: 'defense',
  },
  '000035': {
    children: {
      '000001': 'oig.www',
      '000002': 'www',
    },
    name: 'dhs',
  },
  '000036': {
    children: {
      '000001': 'www',
    },
    name: 'dietaryguidelines',
  },
  '000037': {
    children: {
      '000001': 'standards',
      '000002': 'public-sans',
      '000003': 'accessibility',
      '000004': 'designsystem',
    },
    name: 'digital',
  },
  '000038': {
    children: {
      '000001': 'www',
    },
    name: 'disasterassistance',
  },
  '000039': {
    children: {
      '000001': 'oig.www',
      '000002': 'ntia.www',
      '000003': 'cldp',
    },
    name: 'doc',
  },
  '00003A': {
    children: {
      '000001': 'www',
    },
    name: 'doi',
  },
  '00003B': {
    children: {
      '000001': 'www',
      '000002': 'oalj.www',
    },
    name: 'dol',
  },
  '00003C': {
    children: {},
    name: 'donaciondeorganos',
  },
  '00003D': {
    children: {
      '000001': 'fhwa.www',
      '000002': 'its.www',
      '000003': 'planning.www',
      '000004': 'fhwa.safety.rspcb',
      '000005': 'seaway.www',
      '000006': 'fmcsa.ai',
      '000007': 'its.standards.www',
      '000008': 'its.pcb.www',
      '000009': 'fhwa.environment.www',
      '00000A': 'volpe.www',
      '00000B': 'maritime.www',
      '00000C': 'phmsa.www',
      '00000D': 'railroads',
      '00000E': 'fmcsa.www',
      '00000F': 'highways',
      '000010': 'transit.www',
    },
    name: 'dot',
  },
  '00003E': {
    children: {
      '000001': 'www',
    },
    name: 'drought',
  },
  '00003F': {
    children: {
      '000001': 'archives',
      '000002': 'teens',
      '000003': 'www',
    },
    name: 'drugabuse',
  },
  '000040': {
    children: {
      '000001': 'www2',
      '000002': 'tech',
      '000003': 'oha',
      '000004': 'lincs',
      '000005': 'ifap',
      '000006': 'blog',
      '000007': 'oese',
      '000008': 'collegescorecard',
      '000009': 'fsapartners',
      '00000A': 'sites',
    },
    name: 'ed',
  },
  '000041': {
    children: {
      '000001': 'www',
    },
    name: 'eda',
  },
  '000042': {
    children: {
      '000001': 'ir',
      '000002': 'www',
    },
    name: 'eia',
  },
  '000043': {
    children: {
      '000001': 'www',
    },
    name: 'eisenhowerlibrary',
  },
  '000044': {
    children: {
      '000001': 'www',
    },
    name: 'energystar',
  },
  '000045': {
    children: {
      '000001': 'www',
      '000002': 'espanol',
    },
    name: 'epa',
  },
  '000046': {
    children: {
      '000001': 'www',
    },
    name: 'evaluation',
  },
  '000047': {
    children: {
      '000001': 'grow',
      '000002': 'www',
    },
    name: 'exim',
  },
  '000048': {
    children: {
      '000001': 'legacy',
    },
    name: 'export',
  },
  '000049': {
    children: {
      '000001': 'www',
    },
    name: 'faa',
  },
  '00004A': {
    children: {
      '000001': 'www',
    },
    name: 'fac',
  },
  '00004B': {
    children: {
      '000001': 'www',
    },
    name: 'farmers',
  },
  '00004C': {
    children: {
      '000001': 'www',
    },
    name: 'fcsm',
  },
  '00004D': {
    children: {
      '000001': 'www',
      '000002': 'accessdata.www',
    },
    name: 'fda',
  },
  '00004E': {
    children: {
      '000001': 'www',
    },
    name: 'fdic',
  },
  '00004F': {
    children: {
      '000001': 'webforms',
      '000002': 'dev',
      '000003': 'www',
    },
    name: 'fec',
  },
  '000050': {
    children: {
      '000001': 'tailored',
      '000002': 'www',
    },
    name: 'fedramp',
  },
  '000051': {
    children: {
      '000001': 'www',
      '000002': 'usfa.www',
    },
    name: 'fema',
  },
  '000052': {
    children: {
      '000001': 'www',
    },
    name: 'fincen',
  },
  '000053': {
    children: {
      '000001': 'www',
    },
    name: 'fishwatch',
  },
  '000054': {
    children: {
      '000001': 'www',
      '000002': 'nfipservices',
      '000003': 'agents',
    },
    name: 'floodsmart',
  },
  '000055': {
    children: {
      '000001': 'www',
    },
    name: 'fmc',
  },
  '000056': {
    children: {
      '000001': 'www',
    },
    name: 'foia',
  },
  '000057': {
    children: {
      '000001': 'www',
    },
    name: 'foodsafety',
  },
  '000058': {
    children: {
      '000001': 'www',
    },
    name: 'fordlibrarymuseum',
  },
  '000059': {
    children: {
      '000001': 'www',
    },
    name: 'fpc',
  },
  '00005A': {
    children: {
      '000001': 'www',
    },
    name: 'frtib',
  },
  '00005B': {
    children: {
      '000001': 'www',
      '000002': 'consumer.www',
      '000003': 'consumer',
    },
    name: 'ftc',
  },
  '00005C': {
    children: {
      '000001': 'www',
    },
    name: 'fws',
  },
  '00005D': {
    children: {
      '000001': 'www',
    },
    name: 'genome',
  },
  '00005E': {
    children: {
      '000001': 'beta',
    },
    name: 'get',
  },
  '00005F': {
    children: {
      '000001': 'health2016',
      '000002': 'nca2014',
      '000003': 'science2017',
      '000004': 'carbon2018',
      '000005': 'nca2018',
    },
    name: 'globalchange',
  },
  '000060': {
    children: {
      '000001': 'www',
    },
    name: 'goes-r',
  },
  '000061': {
    children: {
      '000001': 'www',
    },
    name: 'govloans',
  },
  '000062': {
    children: {
      '000001': 'www',
    },
    name: 'gps',
  },
  '000063': {
    children: {
      '000001': 'test',
      '000002': 'staging',
      '000003': 'training',
    },
    name: 'grants',
  },
  '000064': {
    children: {
      '000001': 'www',
      '000002': 'identityequitystudy',
      '000003': 'recycling',
      '000004': 'fedsim',
      '000005': 'tts',
      '000006': 'aas',
      '000007': '10x',
      '000008': 'smartpay.demo',
      '000009': 'cic',
      '00000A': 'digitalcorps',
      '00000B': 'ussm',
      '00000C': 'tts.handbook',
      '00000D': 'smartpay',
      '00000E': 'itvmo',
      '00000F': '18f',
      '000010': 'oes',
    },
    name: 'gsa',
  },
  '000065': {
    children: {
      '000001': 'www',
    },
    name: 'healthcare',
  },
  '000066': {
    children: {
      '000001': 'www',
    },
    name: 'healthit',
  },
  '000067': {
    children: {
      '000001': 'www',
    },
    name: 'helpwithmybank',
  },
  '000068': {
    children: {
      '000001': 'oig',
      '000002': 'betobaccofree',
      '000003': 'betobaccofree.therealcost',
      '000004': 'empowerprogram',
      '000005': 'telehealth',
      '000006': 'asprtracie.files',
      '000007': 'acf.ncsacw',
      '000008': 'asprtracie',
    },
    name: 'hhs',
  },
  '000069': {
    children: {
      '000001': 'historyhub',
    },
    name: 'history',
  },
  '00006A': {
    children: {
      '000001': 'www',
    },
    name: 'hiv',
  },
  '00006B': {
    children: {
      '000001': 'www',
    },
    name: 'hive',
  },
  '00006C': {
    children: {
      '000001': 'poisonhelp',
      '000002': 'bloodstemcell',
      '000003': 'newbornscreening',
      '000004': 'nhsc',
      '000005': 'npdb',
      '000006': 'npdb.www',
      '000007': 'bhw',
      '000008': 'ryanwhite',
      '000009': 'mchb',
      '00000A': 'hab',
      '00000B': 'bphc',
      '00000C': 'www',
    },
    name: 'hrsa',
  },
  '00006D': {
    children: {
      '000001': 'www',
    },
    name: 'hud',
  },
  '00006E': {
    children: {
      '000001': 'www',
    },
    name: 'ice',
  },
  '00006F': {
    children: {
      '000001': 'devicepki',
      '000002': 'www',
    },
    name: 'idmanagement',
  },
  '000070': {
    children: {
      '000001': 'www',
    },
    name: 'ihs',
  },
  '000071': {
    children: {
      '000001': 'www',
    },
    name: 'imls',
  },
  '000072': {
    children: {
      '000001': 'www',
    },
    name: 'invasivespeciesinfo',
  },
  '000073': {
    children: {
      '000001': 'www',
    },
    name: 'irs',
  },
  '000074': {
    children: {
      '000001': 'www',
    },
    name: 'irsauctions',
  },
  '000075': {
    children: {
      '000001': 'www',
    },
    name: 'itap',
  },
  '000076': {
    children: {
      '000001': 'www',
    },
    name: 'jimmycarterlibrary',
  },
  '000077': {
    children: {
      '000001': 'www',
      '000002': 'civilrights',
      '000003': 'oig',
    },
    name: 'justice',
  },
  '000078': {
    children: {
      '000001': 'crd',
    },
    name: 'lbl',
  },
  '000079': {
    children: {
      '000001': 'www',
    },
    name: 'lep',
  },
  '00007A': {
    children: {
      '000001': 'developers',
    },
    name: 'login',
  },
  '00007B': {
    children: {
      '000001': 'www',
    },
    name: 'makinghomeaffordable',
  },
  '00007C': {
    children: {
      '000001': 'www',
    },
    name: 'mbda',
  },
  '00007D': {
    children: {
      '000001': 'www',
    },
    name: 'mcc',
  },
  '00007E': {
    children: {
      '000001': 'www',
    },
    name: 'medicaid',
  },
  '00007F': {
    children: {
      '000001': 'es',
      '000002': 'www',
    },
    name: 'medicare',
  },
  '000080': {
    children: {
      '000001': 'magazine',
    },
    name: 'medlineplus',
  },
  '000081': {
    children: {
      '000001': 'espanol',
      '000002': 'www',
    },
    name: 'mentalhealth',
  },
  '000082': {
    children: {
      '000001': 'www',
    },
    name: 'moneyfactory',
  },
  '000083': {
    children: {
      '000001': 'www',
      '000002': 'arlweb',
    },
    name: 'msha',
  },
  '000084': {
    children: {
      '000001': 'www',
    },
    name: 'mspb',
  },
  '000085': {
    children: {
      '000001': 'www',
    },
    name: 'mymoney',
  },
  '000086': {
    children: {
      '000001': 'www',
    },
    name: 'myplate',
  },
  '000087': {
    children: {
      '000001': 'blogs',
      '000002': 'www-staging',
      '000003': 'beta',
      '000004': 'science',
      '000005': 'www',
      '000006': 'jpl.www',
      '000007': 'solarsystem',
      '000008': 'nsstc.ghrc',
      '000009': 'arc.c3rs',
      '00000A': 'science3',
      '00000B': 'gsfc.etd',
      '00000C': 'earthdata.uat.cdn',
      '00000D': 'earthdata.sit.cdn',
      '00000E': 'earthdata.cdn',
      '00000F': 'ciencia',
      '000010': 'plus',
      '000011': 'www3',
      '000012': 'jpl.podaac',
      '000013': 'larc.eosweb',
      '000014': 'climate',
      '000015': 'jpl.photojournal',
      '000016': 'gsfc.spdf',
      '000017': 'jsc.historycollection',
      '000018': 'earthdata',
      '000019': 'apod',
      '00001A': 'appel',
      '00001B': 'spaceflight',
      '00001C': 'history',
      '00001D': 'earthdata.cmr',
      '00001E': 'mars',
      '00001F': 'science.beta',
      '000020': 'earthobservatory',
    },
    name: 'nasa',
  },
  '000088': {
    children: {
      '000001': 'beta',
      '000002': 'www',
    },
    name: 'ncd',
  },
  '000089': {
    children: {
      '000001': 'www',
    },
    name: 'ncjrs',
  },
  '00008A': {
    children: {
      '000001': 'docs',
      '000002': 'www',
    },
    name: 'nersc',
  },
  '00008B': {
    children: {
      '000001': 'www',
    },
    name: 'nhtsa',
  },
  '00008C': {
    children: {
      '000001': 'www',
    },
    name: 'niem',
  },
  '00008D': {
    children: {
      '000001': 'nhlbi.healthyeating',
      '000002': 'obesityresearch.www',
      '000003': 'downsyndrome',
      '000004': 'niaaa.rethinkingdrinking.www',
      '000005': 'researchtraining',
      '000006': 'niaaa.spectrum.www',
      '000007': 'sharing',
      '000008': 'nichd.safetosleep',
      '000009': 'lrp.www',
      '00000A': 'nhlbi.grasp',
      '00000B': 'nhlbi.catalog',
      '00000C': 'niddk.spin',
      '00000D': 'oir.oacu',
      '00000E': 'olaw',
      '00000F': 'covid19community',
      '000010': 'heal',
      '000011': 'niaaa.www',
      '000012': 'era',
      '000013': 'diversity',
      '000014': 'newsinhealth',
      '000015': 'oir',
      '000016': 'csr.public',
      '000017': 'allofus',
      '000018': 'od.orwh',
      '000019': 'ninr.www',
      '00001A': 'era.www',
      '00001B': 'nida',
      '00001C': 'nichd.espanol',
      '00001D': 'nibib.www',
      '00001E': 'researchfestival',
      '00001F': 'od.ods',
      '000020': 'nccih.www',
      '000021': 'fic.www',
      '000022': 'ncats',
      '000023': 'nidcr.www',
      '000024': 'nidcd.www',
      '000025': 'cc.www',
      '000026': 'nei.www',
      '000027': 'nimhd.www',
      '000028': 'cc',
      '000029': 'nichd.www',
      '00002A': 'nihrecord',
      '00002B': 'niams.www',
      '00002C': 'niddk.www',
      '00002D': 'www',
      '00002E': 'niehs.www',
      '00002F': 'irp',
      '000030': 'hr.directorsawards',
      '000031': 'nimh.www',
      '000032': 'nhlbi.www',
    },
    name: 'nih',
  },
  '00008E': {
    children: {
      '000001': 'www',
      '000002': 'shop',
      '000003': 'itl.www',
      '000004': 'chemdata',
      '000005': 'materialsdata',
      '000006': 'bigdatawg',
      '000007': 'pages',
      '000008': 'csrc',
    },
    name: 'nist',
  },
  '00008F': {
    children: {
      '000001': 'www',
    },
    name: 'nixonlibrary',
  },
  '000090': {
    children: {
      '000002': 'www',
      '000003': 'news',
    },
    name: 'nnlm',
  },
  '000091': {
    children: {
      '000001': 'fisheries.www',
      '000002': 'glerl.coastwatch',
      '000003': 'library.repository',
      '000004': 'afsc.access',
      '000005': 'nwave.carto',
      '000006': 'spc.www',
      '000007': 'nowcoast',
      '000008': 'cameochemicals',
      '000009': 'cameo',
      '00000A': 'nssl.ciflow',
      '00000B': 'nssl.inside',
      '00000C': 'aoml.cwcaribbean',
      '00000D': 'charts',
      '00000E': 'research.qosap',
      '00000F': 'ncep.bldr.madis-data',
      '000010': 'vlab',
      '000011': 'historicalcharts',
      '000012': 'iuufishing',
      '000013': 'iocm',
      '000014': 'marinedebris.clearinghouse',
      '000015': 'nosc',
      '000016': 'goes.www',
      '000017': 'marinedebris.blog',
      '000018': 'coralreef',
      '000019': 'restoration.response.blog',
      '00001A': 'ci',
      '00001B': 'pmel.www',
      '00001C': 'nodc.www',
      '00001D': 'ngdc.www',
      '00001E': 'gsl',
      '00001F': 'arctic',
      '000020': 'coastwatch.eastcoast',
      '000021': 'marinedebris',
      '000022': 'sos',
      '000023': 'nauticalcharts',
      '000024': 'aoml.www',
      '000025': 'coastwatch',
      '000026': 'omao.www',
      '000027': 'ncei.www',
      '000028': 'pifsc.oceanwatch',
      '000029': 'geodesy',
      '00002A': 'psl',
      '00002B': 'nesdis.www',
      '00002C': 'cpo',
      '00002D': 'amdar',
      '00002E': 'climate',
      '00002F': 'coast',
      '000030': 'www',
      '000031': 'oceanexplorer',
      '000032': 'ncdc.www',
      '000033': 'pfeg.coastwatch',
    },
    name: 'noaa',
  },
  '000092': {
    children: {
      '000001': 'www',
    },
    name: 'nps',
  },
  '000093': {
    children: {
      '000001': 'www',
    },
    name: 'nrc',
  },
  '000094': {
    children: {
      '000001': 'www',
    },
    name: 'nro',
  },
  '000095': {
    children: {
      '000001': 'www',
    },
    name: 'nsa',
  },
  '000096': {
    children: {
      '000001': 'www',
      '000002': 'seedfund',
      '000003': 'iucrc',
      '000004': 'beta',
    },
    name: 'nsf',
  },
  '000097': {
    children: {
      '000001': 'its',
      '000002': 'www',
    },
    name: 'ntia',
  },
  '000098': {
    children: {
      '000001': 'www',
    },
    name: 'nutrition',
  },
  '000099': {
    children: {
      '000001': 'www',
    },
    name: 'nwbc',
  },
  '00009A': {
    children: {
      '000001': 'www',
    },
    name: 'obamalibrary',
  },
  '00009B': {
    children: {
      '000001': 'careers',
      '000002': 'www',
    },
    name: 'occ',
  },
  '00009C': {
    children: {
      '000001': 'www',
    },
    name: 'ofia',
  },
  '00009D': {
    children: {
      '000001': 'extapps2',
      '000002': 'www',
    },
    name: 'oge',
  },
  '00009E': {
    children: {
      '000001': 'www',
    },
    name: 'ojjdp',
  },
  '00009F': {
    children: {
      '000001': 'bja',
      '000002': 'ojjdp',
      '000003': 'nij',
    },
    name: 'ojp',
  },
  '0000A0': {
    children: {
      '000001': 'www',
    },
    name: 'onhir',
  },
  '0000A1': {
    children: {
      '000002': 'www',
    },
    name: 'onrr',
  },
  '0000A2': {
    children: {
      '000001': 'www',
    },
    name: 'opm',
  },
  '0000A3': {
    children: {
      '000001': 'orise',
    },
    name: 'orau',
  },
  '0000A4': {
    children: {
      '000001': 'www',
    },
    name: 'organdonor',
  },
  '0000A5': {
    children: {
      '000001': 'carve',
      '000002': 'modis',
      '000003': 'daac-news',
      '000004': 'daac',
    },
    name: 'ornl',
  },
  '0000A6': {
    children: {
      '000001': 'www',
    },
    name: 'ourdocuments',
  },
  '0000A7': {
    children: {
      '000001': 'www',
      '000002': 'abilityone',
    },
    name: 'oversight',
  },
  '0000A8': {
    children: {
      '000001': 'www',
    },
    name: 'pbrb',
  },
  '0000A9': {
    children: {
      '000001': 'www',
    },
    name: 'performance',
  },
  '0000AA': {
    children: {
      '000001': 'www',
    },
    name: 'pic',
  },
  '0000AB': {
    children: {},
    name: 'plainlanguage',
  },
  '0000AC': {
    children: {},
    name: 'presidentialinnovationfellows',
  },
  '0000AD': {
    children: {
      '000001': 'www',
    },
    name: 'ready',
  },
  '0000AE': {
    children: {
      '000001': 'www',
    },
    name: 'reaganlibrary',
  },
  '0000AF': {
    children: {
      '000001': 'blog',
      '000002': 'ncsacw',
      '000003': 'store',
      '000004': 'www',
    },
    name: 'samhsa',
  },
  '0000B0': {
    children: {
      '000001': 'www',
    },
    name: 'sandia',
  },
  '0000B1': {
    children: {
      '000001': 'www',
    },
    name: 'schoolsafety',
  },
  '0000B2': {
    children: {
      '000000': '',
    },
    name: 'search',
  },
  '0000B3': {
    children: {
      '000001': 'careers',
      '000002': 'www',
    },
    name: 'secretservice',
  },
  '0000B4': {
    children: {
      '000001': 'www',
    },
    name: 'section508',
  },
  '0000B5': {
    children: {
      '000001': 'help.www',
    },
    name: 'senate',
  },
  '0000B6': {
    children: {},
    name: 'sftool',
  },
  '0000B7': {
    children: {
      '000001': 'wwwtest',
      '000002': 'www',
      '000003': 'blog',
      '000004': 'faq',
      '000005': 'oig-files',
      '000006': 'oig-demo',
      '000007': 'oig',
    },
    name: 'ssa',
  },
  '0000B8': {
    children: {
      '000001': 'www',
    },
    name: 'ssab',
  },
  '0000B9': {
    children: {
      '000001': 'www',
    },
    name: 'sss',
  },
  '0000BA': {
    children: {
      '000001': 'www',
      '000002': 'palestinianaffairs',
      '000003': 'ylai',
      '000004': 'yali',
      '000005': 'statemag',
      '000006': 'careers',
    },
    name: 'state',
  },
  '0000BB': {
    children: {
      '000001': 'www',
    },
    name: 'statspolicy',
  },
  '0000BC': {
    children: {
      '000001': 'prod',
      '000002': 'www',
    },
    name: 'stb',
  },
  '0000BD': {
    children: {
      '000001': 'www',
    },
    name: 'stopalcoholabuse',
  },
  '0000BE': {
    children: {
      '000001': 'www',
    },
    name: 'stopbullying',
  },
  '0000BF': {
    children: {
      '000001': 'www',
    },
    name: 'tigta',
  },
  '0000C0': {
    children: {
      '000001': 'www7',
      '000002': 'www',
    },
    name: 'transportation',
  },
  '0000C1': {
    children: {
      '000001': 'www',
      '000002': 'fiscal.tcvs',
      '000003': 'oig',
      '000004': 'ofac',
      '000005': 'fiscal',
      '000006': 'home',
    },
    name: 'treasury',
  },
  '0000C2': {
    children: {
      '000002': 'www',
    },
    name: 'treasurydirect',
  },
  '0000C3': {
    children: {
      '000001': 'www',
    },
    name: 'tsa',
  },
  '0000C4': {
    children: {
      '000001': 'www',
    },
    name: 'ttb',
  },
  '0000C5': {
    children: {
      '000001': 'niccs',
    },
    name: 'us-cert',
  },
  '0000C6': {
    children: {
      '000001': 'benefits-tool',
      '000002': 'blog',
      '000003': 'www',
    },
    name: 'usa',
  },
  '0000C7': {
    children: {
      '000001': 'www',
    },
    name: 'usability',
  },
  '0000C8': {
    children: {
      '000001': 'www',
    },
    name: 'usagm',
  },
  '0000C9': {
    children: {
      '000001': '2012-2017',
      '000002': 'oig',
      '000003': 'blog',
      '000004': 'www',
    },
    name: 'usaid',
  },
  '0000CA': {
    children: {
      '000001': 'www',
    },
    name: 'usbr',
  },
  '0000CB': {
    children: {
      '000001': 'www',
    },
    name: 'uscc',
  },
  '0000CC': {
    children: {
      '000001': 'www',
      '000002': 'my',
      '000003': 'preview',
    },
    name: 'uscis',
  },
  '0000CD': {
    children: {
      '000001': 'bm',
      '000002': 'hk',
    },
    name: 'usconsulate',
  },
  '0000CE': {
    children: {
      '000001': 'ca7.media',
      '000002': 'cvb.www',
      '000003': 'cvb.www3',
      '000004': 'ilsp.www',
      '000005': 'arep.www',
      '000006': 'mssp.www',
      '000007': 'mow.www',
      '000008': 'msnb.www',
      '000009': 'gand.www',
      '00000A': 'ncwba.www',
      '00000B': 'pacer',
      '00000C': 'flnb.www',
      '00000D': 'njd.www',
      '00000E': 'lawb.www',
      '00000F': 'nep.www',
      '000010': 'rid.www',
      '000011': 'mssd.www',
      '000012': 'ilsd.www',
      '000013': 'oknb.www',
      '000014': 'jpml.www',
      '000015': 'caep.www',
      '000016': 'mad.www',
      '000017': 'msnd.www',
      '000018': 'ca7.www',
      '000019': 'wiwb.www',
      '00001A': 'moeb.www',
      '00001B': 'are.www',
      '00001C': 'ned.www',
      '00001D': 'tneb.www',
      '00001E': 'deb.www',
      '00001F': 'pawd.www',
      '000020': 'scb.www',
      '000021': 'mssd.ecf',
      '000022': 'cit.www',
      '000023': 'wvsd.www',
      '000024': 'flsd.www',
      '000025': 'cacb.www',
      '000026': 'utb.www',
      '000027': 'www',
    },
    name: 'uscourts',
  },
  '0000CF': {
    children: {
      '000001': 'nass.www',
      '000002': 'nal.i5k',
      '000003': 'ams.search',
      '000004': 'egov.sc.wcc',
      '000005': 'fns.farmtoschoolcensus',
      '000006': 'fns.snaptoskills',
      '000007': 'fns.wicbreastfeeding',
      '000008': 'ars.aglab',
      '000009': 'scinet',
      '00000A': 'fns.professionalstandards',
      '00000B': 'nesr',
      '00000C': 'fns.wicworks',
      '00000D': 'fns.snaped',
      '00000E': 'nfc',
      '00000F': 'rma',
      '000010': 'nrcs.www',
      '000011': 'aphis.www',
      '000012': 'nal.www',
      '000013': 'ers.www',
      '000014': 'nfc.help',
      '000015': 'fns.www',
    },
    name: 'usda',
  },
  '0000D0': {
    children: {
      '000001': 'cops',
    },
    name: 'usdoj',
  },
  '0000D1': {
    children: {
      '000001': 'mv',
      '000002': 'bi',
      '000003': 'pw',
      '000004': 'sample',
      '000005': 'fm',
      '000006': 'om',
      '000007': 'so',
      '000008': 'tg',
      '000009': 'sl',
      '00000A': 'sd',
      '00000B': 'nl',
      '00000C': 'gm',
      '00000D': 'ls',
      '00000E': 'sz',
      '00000F': 'bz',
      '000010': 'gq',
      '000011': 'se',
      '000012': 'sa',
      '000013': 'mw',
      '000014': 'zw',
      '000015': 'km',
      '000016': 'kw',
      '000017': 'cu',
      '000018': 'no',
      '000019': 'bw',
      '00001A': 'bs',
      '00001B': 'mt',
      '00001C': 'jm',
      '00001D': 'ie',
      '00001E': 'lr',
      '00001F': 'la',
      '000020': 'ws',
      '000021': 'ca',
      '000022': 'bh',
      '000023': 'sy',
      '000024': 'fi',
      '000025': 'pg',
      '000026': 'cg',
      '000027': 'lk',
      '000028': 'ye',
      '000029': 'hr',
      '00002A': 'pt',
      '00002B': 'ly',
      '00002C': 'fj',
      '00002D': 'dk',
      '00002E': 'ec',
      '00002F': 'lu',
      '000030': 'bf',
      '000031': 'np',
      '000032': 'sn',
      '000033': 'uy',
      '000034': 'cy',
      '000035': 'be',
      '000036': 'rs',
      '000037': 'au',
      '000038': 'si',
      '000039': 'zm',
      '00003A': 'ir',
      '00003B': 've',
      '00003C': 'ne',
      '00003D': 'gn',
      '00003E': 'ni',
      '00003F': 'va',
      '000040': 'al',
      '000041': 'tr',
      '000042': 'ke',
      '000043': 'bb',
      '000044': 'qa',
      '000045': 'iq',
      '000046': 'lb',
      '000047': 'cr',
      '000048': 'az',
      '000049': 'kg',
      '00004A': 'ae',
      '00004B': 'at',
      '00004C': 'za',
      '00004D': 'td',
      '00004E': 'id',
      '00004F': 'nz',
      '000050': 'sg',
      '000051': 'fr',
      '000052': 'sample2',
      '000053': 'hu',
      '000054': 'tt',
      '000055': 'mg',
      '000056': 'mz',
      '000057': 'bo',
      '000058': 'cm',
      '000059': 'ar',
      '00005A': 'vn',
      '00005B': 'pk',
      '00005C': 'sv',
      '00005D': 'bd',
      '00005E': 'pa',
      '00005F': 'jp',
      '000060': 'ci',
      '000061': 'pe',
      '000062': 'kz',
      '000063': 'my',
      '000064': 'de',
      '000065': 'es',
      '000066': 'co',
      '000067': 'ph',
      '000068': 'ng',
      '000069': 'gt',
      '00006A': 'af',
      '00006B': 'kr',
      '00006C': 'mx',
      '00006D': 'ge',
      '00006E': 'uk',
      '00006F': 'br',
    },
    name: 'usembassy',
  },
  '0000D2': {
    children: {
      '000001': 'www',
      '000002': 'lpdaac',
      '000003': 'umesc',
    },
    name: 'usgs',
  },
  '0000D3': {
    children: {
      '000001': 'www',
    },
    name: 'usmarshals',
  },
  '0000D4': {
    children: {
      '000001': 'usoecd',
      '000002': 'nato',
      '000003': 'useu',
      '000004': 'usunrome',
      '000005': 'asean',
      '000006': 'vienna',
      '000007': 'usun',
      '000008': 'geneva',
    },
    name: 'usmission',
  },
  '0000D5': {
    children: {
      '000001': 'www',
    },
    name: 'uspsoig',
  },
  '0000D6': {
    children: {
      '000001': 'www',
      '000002': '10millionpatents',
      '000003': 'foiadocuments',
    },
    name: 'uspto',
  },
  '0000D7': {
    children: {
      '000001': 'www',
      '000002': 'cem.gravelocator',
      '000003': 'vba.vaonce',
      '000004': 'gibill.www',
      '000005': 'pay.www',
      '000006': 'vba.inquiry',
      '000007': 'section508.www',
      '000008': 'vrss',
      '000009': 'rcv.www',
      '00000A': 'research.cerc.www',
      '00000B': 'research.cider.www',
      '00000C': 'oedca.www',
      '00000D': 'sci.www',
      '00000E': 'vetcenter.www',
      '00000F': 'research.brrc.www',
      '000010': 'ehrm.www',
      '000011': 'fss.www',
      '000012': 'research.virec.www',
      '000013': 'sep.www',
      '000014': 'innovation.www',
      '000015': 'visn15.www',
      '000016': 'research.cshiip.www',
      '000017': 'osp.www',
      '000018': 'visn19.www',
      '000019': 'choose',
      '00001A': 'psychologytraining.www',
      '00001B': 'oefoif.www',
      '00001C': 'southeast.www',
      '00001D': 'visn23.www',
      '00001E': 'ebenefits.www',
      '00001F': 'fsc.www',
      '000020': 'med.iowa-city.research.www',
      '000021': 'research.chic.www',
      '000022': 'research.portlandcoin.www',
      '000023': 'research.vision.www',
      '000024': 'research.cmc3.www',
      '000025': 'visn10.www',
      '000026': 'research.avreap.www',
      '000027': 'visn16.www',
      '000028': 'research.seattledenvercoin.www',
      '000029': 'heartoftexas.www',
      '00002A': 'socialwork.www',
      '00002B': 'research.peprec.www',
      '00002C': 'visn6.www',
      '00002D': 'visn9.www',
      '00002E': 'energy.www',
      '00002F': 'research.ccdor.www',
      '000030': 'research.eric.seattle.www',
      '000031': 'poplarbluff.www',
      '000032': 'research.hsrd.durham.www',
      '000033': 'visn12.www',
      '000034': 'med.visn20.www',
      '000035': 'research.cadre.www',
      '000036': 'visn21.www',
      '000037': 'centralwesternmass.www',
      '000038': 'northport.www',
      '000039': 'desertpacific.www',
      '00003A': 'leavenworth.www',
      '00003B': 'valu.www',
      '00003C': 'research.annarbor.www',
      '00003D': 'elpaso.www',
      '00003E': 'dublin.www',
      '00003F': 'wichita.www',
      '000040': 'sheridan.www',
      '000041': 'epilepsy.www',
      '000042': 'marion.www',
      '000043': 'roseburg.www',
      '000044': 'jackson.www',
      '000045': 'texasvalley.www',
      '000046': 'hampton.www',
      '000047': 'bva.www',
      '000048': 'tomah.www',
      '000049': 'kansascity.www',
      '00004A': 'research.ci2i.www',
      '00004B': 'bigspring.www',
      '00004C': 'ironmountain.www',
      '00004D': 'southernoregon.www',
      '00004E': 'research.hsrd.annarbor.www',
      '00004F': 'salem.www',
      '000050': 'wallawalla.www',
      '000051': 'visn8.www',
      '000052': 'research.amputation.www',
      '000053': 'columbus.www',
      '000054': 'vacareers.www',
      '000055': 'huntington.www',
      '000056': 'research.cindrr.www',
      '000057': 'research.hsrd.houston.www',
      '000058': 'maine.www',
      '000059': 'montana.www',
      '00005A': 'dieteticinternship.www',
      '00005B': 'patientcare.www',
      '00005C': 'grandjunction.www',
      '00005D': 'research.vacsp.www',
      '00005E': 'chicago.www',
      '00005F': 'coatesville.www',
      '000060': 'fayettevillear.www',
      '000061': 'hawaii.www',
      '000062': 'mountainhome.www',
      '000063': 'wilkes-barre.www',
      '000064': 'westpalmbeach.www',
      '000065': 'birmingham.www',
      '000066': 'cheyenne.www',
      '000067': 'polytrauma.www',
      '000068': 'veterantraining.www',
      '000069': 'whiteriver.www',
      '00006A': 'developer',
      '00006B': 'northernindiana.www',
      '00006C': 'newjersey.www',
      '00006D': 'caregiver.www',
      '00006E': 'centraliowa.www',
      '00006F': 'research.ncrar.www',
      '000070': 'reno.www',
      '000071': 'clarksburg.www',
      '000072': 'danville.www',
      '000073': 'topeka.www',
      '000074': 'boise.www',
      '000075': 'iowacity.www',
      '000076': 'shreveport.www',
      '000077': 'memphis.www',
      '000078': 'miami.www',
      '000079': 'fresno.www',
      '00007A': 'manchester.www',
      '00007B': 'hudsonvalley.www',
      '00007C': 'vaforvets.www',
      '00007D': 'philadelphia.www',
      '00007E': 'bronx.www',
      '00007F': 'lomalinda.www',
      '000080': 'fargo.www',
      '000081': 'columbiamo.www',
      '000082': 'salisbury.www',
      '000083': 'lexington.www',
      '000084': 'wilmington.www',
      '000085': 'acquisitionacademy.www',
      '000086': 'longbeach.www',
      '000087': 'oit.ea.www',
      '000088': 'detroit.www',
      '000089': 'healthquality.www',
      '00008A': 'nutrition.www',
      '00008B': 'centralalabama.www',
      '00008C': 'visn4.www',
      '00008D': 'prosthetics.www',
      '00008E': 'siouxfalls.www',
      '00008F': 'spokane.www',
      '000090': 'houston.www',
      '000091': 'cleveland.www',
      '000092': 'caribbean.www',
      '000093': 'volunteer.www',
      '000094': 'newengland.www',
      '000095': 'tuscaloosa.www',
      '000096': 'biloxi.www',
      '000097': 'chillicothe.www',
      '000098': 'louisville.www',
      '000099': 'denver.www',
      '00009A': 'orlando.www',
      '00009B': 'oklahoma.www',
      '00009C': 'cincinnati.www',
      '00009D': 'ruralhealth.www',
      '00009E': 'simlearn.www',
      '00009F': 'littlerock.www',
      '0000A0': 'fayettevillenc.www',
      '0000A1': 'research.choir.www',
      '0000A2': 'tampa.www',
      '0000A3': 'research.cherp.www',
      '0000A4': 'connecticut.www',
      '0000A5': 'indianapolis.www',
      '0000A6': 'lasvegas.www',
      '0000A7': 'pugetsound.www',
      '0000A8': 'madison.www',
      '0000A9': 'saltlakecity.www',
      '0000AA': 'columbiasc.www',
      '0000AB': 'syracuse.www',
      '0000AC': 'sandiego.www',
      '0000AD': 'tucson.www',
      '0000AE': 'losangeles.www',
      '0000AF': 'digital',
      '0000B0': 'research.aptcenter.www',
      '0000B1': 'parkinsons.www',
      '0000B2': 'nyharbor.www',
      '0000B3': 'tennesseevalley.www',
      '0000B4': 'northtexas.www',
      '0000B5': 'stlouis.www',
      '0000B6': 'oprm.www',
      '0000B7': 'sanfrancisco.www',
      '0000B8': 'butler.www',
      '0000B9': 'saginaw.www',
      '0000BA': 'connectedcare',
      '0000BB': 'centraltexas.www',
      '0000BC': 'richmond.www',
      '0000BD': 'hines.www',
      '0000BE': 'phoenix.www',
      '0000BF': 'fhcc.lovell.www',
      '0000C0': 'patientsafety.www',
      '0000C1': 'blackhills.www',
      '0000C2': 'martinsburg.www',
      '0000C3': 'visn2.www',
      '0000C4': 'neworleans.www',
      '0000C5': 'alexandria.www',
      '0000C6': 'erie.www',
      '0000C7': 'stcloud.www',
      '0000C8': 'bedford.www',
      '0000C9': 'prevention.www',
      '0000CA': 'beckley.www',
      '0000CB': 'warrelatedillness.www',
      '0000CC': 'amarillo.www',
      '0000CD': 'womenshealth.www',
      '0000CE': 'altoona.www',
      '0000CF': 'diversity.www',
      '0000D0': 'milwaukee.www',
      '0000D1': 'providence.www',
      '0000D2': 'durham.www',
      '0000D3': 'charleston.www',
      '0000D4': 'lebanon.www',
      '0000D5': 'albany.www',
      '0000D6': 'asheville.www',
      '0000D7': 'bath.www',
      '0000D8': 'research.herc.www',
      '0000D9': 'accesstocare.www',
      '0000DA': 'boston.www',
      '0000DB': 'discover',
      '0000DC': 'northerncalifornia.www',
      '0000DD': 'alaska.www',
      '0000DE': 'move.www',
      '0000DF': 'southtexas.www',
      '0000E0': 'maryland.www',
      '0000E1': 'mobile',
      '0000E2': 'augusta.www',
      '0000E3': 'northflorida.www',
      '0000E4': 'minneapolis.www',
      '0000E5': 'buffalo.www',
      '0000E6': 'dayton.www',
      '0000E7': 'paloalto.www',
      '0000E8': 'atlanta.www',
      '0000E9': 'portland.www',
      '0000EA': 'muskogee.www',
      '0000EB': 'mentalhealth.www',
      '0000EC': 'ethics.www',
      '0000ED': 'nebraska.www',
      '0000EE': 'annarbor.www',
      '0000EF': 'publichealth.www',
      '0000F0': 'washingtondc.www',
      '0000F1': 'hepatitis.www',
      '0000F2': 'pbm.www',
      '0000F3': 'hiv.www',
      '0000F4': 'research.queri.www',
      '0000F5': 'battlecreek.www',
      '0000F6': 'custhelp.iris',
      '0000F7': 'albuquerque.www',
      '0000F8': 'cem.www',
      '0000F9': 'ptsd.www',
      '0000FA': 'department',
      '0000FB': 'cfm.www',
      '0000FC': 'mirecc.www',
      '0000FD': 'myhealth.www',
      '0000FE': 'baypines.www',
      '0000FF': 'research.www',
      '000100': 'veteranshealthlibrary.www',
      '000101': 'news',
      '000102': 'blogs.www',
      '000103': 'oit.www',
      '000104': 'benefits',
      '000105': 'research.hsrd.www',
    },
    name: 'va',
  },
  '0000D8': {
    children: {
      '000001': 'www',
    },
    name: 'vaoig',
  },
  '0000D9': {
    children: {
      '000001': 'www',
    },
    name: 'vcf',
  },
  '0000DA': {
    children: {},
    name: 'vote',
  },
  '0000DB': {
    children: {
      '000001': 'www',
    },
    name: 'worker',
  },
  '0000DC': {
    children: {
      '000001': 'www',
    },
    name: 'wwtg',
  },
};

assert util.getDomain('search', domains) == '0000B2';
assert util.getDomain('cloud', domains) == '000024';
assert util.getDomain('vote', domains) == '0000DA';

{
  domains: domains,
}
