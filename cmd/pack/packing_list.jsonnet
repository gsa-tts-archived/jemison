[
  {
    # `pack` will produce this file
    host: "tts.gov",
    # These will be paths into the `extract` bucket
    # that `pack` will use
    constituents: [
      "tts.gsa.gov",
      "18f.gsa.gov",
      "search.gov",
      "cloud.gov"
    ]
  },
  # In this case, the file and consistuents are the same.
  # Although wordy, this uses the same structure for all builds.
  {
    host: "grants.gov",
    constituents: [
      "grants.gov"
    ],
  },
]