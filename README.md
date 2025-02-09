# Cross Ledger Transaction Consistency for Financial Auditing

This repository contains the reference implementation of the protocols presented in our paper **"Cross Ledger
Transaction Consistency for Financial Auditing"**.
The work introduces two protocols, **CLOSC** and **CLOLC**, designed to enable auditors and regulatory committees to
verify that transactions across different ledgers are consistent, while preserving privacy and scalability in
large-scale financial networks.

## Paper

Our paper is published in the Advances in Financial Technologies (AFT) 2024 proceedings.
To cite our work, please use the following BibTeX entry:

```bibtex
@InProceedings{koutsos_et_al:LIPIcs.AFT.2024.4,
  author =	{Koutsos, Vlasis and Tian, Xiangan and Papadopoulos, Dimitrios and Chatzopoulos, Dimitris},
  title =	{{Cross Ledger Transaction Consistency for Financial Auditing}},
  booktitle =	{6th Conference on Advances in Financial Technologies (AFT 2024)},
  pages =	{4:1--4:25},
  series =	{Leibniz International Proceedings in Informatics (LIPIcs)},
  ISBN =	{978-3-95977-345-4},
  ISSN =	{1868-8969},
  year =	{2024},
  volume =	{316},
  editor =	{B\"{o}hme, Rainer and Kiffer, Lucianna},
  publisher =	{Schloss Dagstuhl -- Leibniz-Zentrum f{\"u}r Informatik},
  address =	{Dagstuhl, Germany},
  URL =		{https://drops.dagstuhl.de/entities/document/10.4230/LIPIcs.AFT.2024.4},
  URN =		{urn:nbn:de:0030-drops-209409},
  doi =		{10.4230/LIPIcs.AFT.2024.4},
  annote =	{Keywords: Financial auditing, Two-tier ledger architecture, Smart contracts, Transaction privacy, Financial entity unlinkability}
}
```

For additional details, see the paper on the [Cryptology ePrint Archive](https://eprint.iacr.org/2024/1155).

## Project Structure

```text
.
├── benchmark
│   ├── clolc           # CLOLC: benchmark contracts, internal modules, and scripts (e.g., scripts for initialization, transaction recording, consistency examination, result verification)
│   ├── closc           # CLOSC: similar structure as clolc for corresponding benchmarks
│   └── timecounter     # Utility for time counting
├── internal            # Core modules (auditor, committee, organization, transaction, crypto, constants) for both CLOLC and CLOSC
├── script              # Setup scripts (e.g., setup.sh)
├── LICENSE
├── README.md
└── go.mod, go.sum
```

## Important Notes on Benchmarking

- Benchmark Execution:
  To run any benchmark, you most first change directory into the corresponding script folder.
  For example:
    - For the **CLOLC** benchmarks, navigate to `benchmark/clolc/script` and then run the desired benchmark script.
    - For the **CLOSC** benchmarks, navigate to `benchmark/closc/script` and then run the appropriate benchmark script.
- Common Benchmark Scripts:
    - `run_all.sh`: Runs all benchmarks for the corresponding protocol.
    - `run_off_chain.sh`: Runs benchmarks for the off-chain phase of the corresponding protocol.
    - `run_on_chain.sh`: Runs benchmarks for the on-chain phase of the corresponding protocol.
- Environment Recommendation:
  For reproducible performance and to fully leverage our benchmarks, we recommend using Linux (Ubuntu) machine.

# Setup

Before running any benchmarks or tests, install the necessary dependencies by running the following script:

```bash
./scripts/setup.sh
```

# Contributing

Contributions, suggestions, and bug reports are welcome! Please open an issue or submit a pull request.

# License

This project is licensed under the MIT License.

# Contact

For further inquiries, please reach out to:

- Vlasis Koutsos - vkoutsos@cse.ust.hk
- Xiangan Tian - xtianae@cse.ust.hk