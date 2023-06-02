import argparse
import yaml


def generate_fablo_config(
    output_filename,
    chaincode_name,
    chaincode_dir,
    num_orderers,
    num_orgs,
    num_auditors,
):
    """Generate Fablo config."""
    config = {
        "$schema": "https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/schema.json",
        "global": {"fabricVersion": "2.4.2", "tls": True},
        "orgs": [],
        "channels": [{"name": "mychannel", "orgs": []}],
        "chaincodes": [
            {
                "name": chaincode_name,
                "version": "0.1.0",
                "lang": "golang",
                "channel": "mychannel",
                "init": '{"Args":[]}',
                "endorsement": "",
                "directory": chaincode_dir,
            }
        ],
    }

    config["orgs"].append(
        {
            "organization": {"name": "Orderer", "domain": "root.com"},
            "orderers": [
                {
                    "groupName": "group1",
                    "prefix": "orderer",
                    "type": "raft",
                    "instances": num_orderers,
                }
            ],
        }
    )

    endorsement_list = []

    for i in range(1, num_orgs + 1):
        org_name = f"Org{i}"
        org_domain = f"org{i}.example.com"
        config["orgs"].append(
            {
                "organization": {"name": org_name, "domain": org_domain},
                "peer": {"instances": 1},
            }
        )
        config["channels"][0]["orgs"].append(
            {"name": org_name, "peers": ["peer0"]}
        )
        endorsement_list.append(f"'{org_name}MSP.member'")

    for i in range(1, num_auditors + 1):
        aud_name = f"Aud{i}"
        aud_domain = f"aud{i}.example.com"
        config["orgs"].append(
            {
                "organization": {"name": aud_name, "domain": aud_domain},
                "peer": {"instances": 1},
            }
        )
        config["channels"][0]["orgs"].append(
            {"name": aud_name, "peers": ["peer0"]}
        )
        endorsement_list.append(f"'{aud_name}MSP.member'")

    config["chaincodes"][0]["endorsement"] = (
        "AND(" + ", ".join(endorsement_list) + ")"
    )

    with open(output_filename, "w") as file:
        yaml.dump(config, file)


def main():
    """Main function to handle command line arguments."""
    parser = argparse.ArgumentParser(description="Generate Fablo config")
    parser.add_argument(
        "--output_filename",
        default="fablo_config.yaml",
        type=str,
        help="Output filename",
    )
    parser.add_argument(
        "--chaincode_name",
        default="auti-local-chain",
        type=str,
        help="Chaincode name",
    )
    parser.add_argument(
        "--chaincode_dir",
        default="contract/clolc_local_chain",
        type=str,
        help="Chaincode directory",
    )
    parser.add_argument(
        "--num_orderers", default=1, type=int, help="Number of orderers"
    )
    parser.add_argument(
        "--num_orgs", default=1, type=int, help="Number of organizations"
    )
    parser.add_argument(
        "--num_auditors", default=1, type=int, help="Number of auditors"
    )
    args = parser.parse_args()

    generate_fablo_config(
        args.output_filename,
        args.chaincode_name,
        args.chaincode_dir,
        args.num_orderers,
        args.num_orgs,
        args.num_auditors,
    )


if __name__ == "__main__":
    main()
