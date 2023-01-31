/*
  Warnings:

  - You are about to drop the column `from` on the `Transactions` table. All the data in the column will be lost.
  - You are about to drop the column `to` on the `Transactions` table. All the data in the column will be lost.
  - You are about to drop the column `walletId` on the `Transactions` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "Transactions" DROP CONSTRAINT "Transactions_walletId_fkey";

-- AlterTable
ALTER TABLE "Transactions" DROP COLUMN "from",
DROP COLUMN "to",
DROP COLUMN "walletId";

-- CreateTable
CREATE TABLE "_TransactionsToWallet" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "_TransactionsToWallet_AB_unique" ON "_TransactionsToWallet"("A", "B");

-- CreateIndex
CREATE INDEX "_TransactionsToWallet_B_index" ON "_TransactionsToWallet"("B");

-- AddForeignKey
ALTER TABLE "_TransactionsToWallet" ADD CONSTRAINT "_TransactionsToWallet_A_fkey" FOREIGN KEY ("A") REFERENCES "Transactions"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_TransactionsToWallet" ADD CONSTRAINT "_TransactionsToWallet_B_fkey" FOREIGN KEY ("B") REFERENCES "Wallet"("id") ON DELETE CASCADE ON UPDATE CASCADE;
