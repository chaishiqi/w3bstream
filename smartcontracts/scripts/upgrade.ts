import { ethers, upgrades } from 'hardhat';

async function main() {

  // if (process.env.W3BSTREAM_TASK_MANAGER) {
  //   const W3bstreamTaskManager = await ethers.getContractFactory('W3bstreamTaskManager');
  //   await upgrades.forceImport(process.env.W3BSTREAM_TASK_MANAGER, W3bstreamTaskManager);
  //   await upgrades.upgradeProxy(process.env.W3BSTREAM_TASK_MANAGER, W3bstreamTaskManager, {});
  //   console.log(`Upgrade W3bstreamTaskManager ${process.env.W3BSTREAM_TASK_MANAGER} successfull!`);
  // }



  if (process.env.W3BSTREAM_TASK_MANAGER) {
    const W3bstreamTaskManager = await ethers.getContractFactory('W3bstreamTaskManager');
    await upgrades.upgradeProxy(process.env.W3BSTREAM_TASK_MANAGER, W3bstreamTaskManager, {});
    console.log(`Upgrade W3bstreamDebits ${process.env.W3BSTREAM_TASK_MANAGER} successfull!`);
  }
}

main().catch(err => {
  console.error(err);
  process.exitCode = 1;
});
