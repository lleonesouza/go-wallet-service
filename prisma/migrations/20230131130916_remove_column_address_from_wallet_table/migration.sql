/*
  Warnings:

  - You are about to drop the column `address` on the `Wallet` table. All the data in the column will be lost.

*/
-- DropIndex
DROP INDEX "Wallet_address_key";

-- AlterTable
ALTER TABLE "Wallet" DROP COLUMN "address";
