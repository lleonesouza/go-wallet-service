/*
  Warnings:

  - You are about to drop the `Shopkeeper` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Transactions` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_TransactionsToWallet` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Shopkeeper" DROP CONSTRAINT "Shopkeeper_walletId_fkey";

-- DropForeignKey
ALTER TABLE "_TransactionsToWallet" DROP CONSTRAINT "_TransactionsToWallet_A_fkey";

-- DropForeignKey
ALTER TABLE "_TransactionsToWallet" DROP CONSTRAINT "_TransactionsToWallet_B_fkey";

-- DropTable
DROP TABLE "Shopkeeper";

-- DropTable
DROP TABLE "Transactions";

-- DropTable
DROP TABLE "_TransactionsToWallet";

-- CreateTable
CREATE TABLE "Transaction" (
    "id" TEXT NOT NULL,
    "value" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "walletId" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Transaction_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Question" (
    "id" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "input" TEXT NOT NULL,
    "output" TEXT,
    "userId" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Question_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "Transaction" ADD CONSTRAINT "Transaction_walletId_fkey" FOREIGN KEY ("walletId") REFERENCES "Wallet"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Question" ADD CONSTRAINT "Question_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
