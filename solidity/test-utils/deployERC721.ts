import { GravityERC721} from "../typechain/GravityERC721";
import { TestERC721A } from "../typechain/TestERC721A";
import { ethers } from "hardhat";
import { Signer } from "ethers";
import { deployContracts } from "./index";

export async function deployContractsERC721(
  gravityId: string = "foo",
  validators: Signer[],
  powers: number[],
) {

  const {
    gravity,
    testERC20
  } = await deployContracts(gravityId, validators, powers);

  const TestERC721 = await ethers.getContractFactory("TestERC721A");
  const testERC721= (await TestERC721.deploy()) as TestERC721A;

  const GravityERC721 = await ethers.getContractFactory("GravityERC721");
  const gravityERC721 = (await GravityERC721.deploy(
    gravity.address
  )) as GravityERC721;

  return { gravity, gravityERC721, testERC721,  testERC20 };
}
