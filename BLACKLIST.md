# Blacklist info

The black list is a newline delimited file of wallet addresses. It can also support comments with the `#` character.

## Default Location

By default, the feechain binary looks for the file `./.fch/blaklist.txt`.

## Example File
```
fee1spshr72utf6rwxseaz339j09ed8p6f8k0vtvpv
fee1uyshu2jgv8w465yc8kkny36thlt2wvel3c7mck  # This is a comment
fee1r4zyyjqrulf935a479sgqlpa78kz7zlc7h826d

```

## Details

Each transaction added to the tx-pool has its `to` and `from` address checked against this blacklist. 
If there is a hit, the transaction is considered invalid and is dropped from the tx-pool.