
## x
* `multisig` - 30542 1

## y

### Pay-to-PublicKey（P2PK）

```bash
锁定脚本：<Public Key> OP_CHECKSIG

解锁脚本：<Signature from Private Key>

验证时的组合脚本：<Signature from Private Key> <Public Key> OP_CHECKSIG
```

### Pay-to-PublicKey Hash（P2PKH）

```bash
锁定脚本：OP_DUP OP_HASH160 <Public KeyHash> OP_EQUAL OP_CHECKSIG

解锁脚本：<Signature> <Public Key>

验证时的组合脚本：<Signature> <Public Key> OP_DUP OP_HASH160 <Public KeyHash> OP_EQUAL OP_CHECKSIG
```

### 多重签名P2MS(Multiple Signatures)

```bash
锁定脚本为

OP_M
<PublicKey 1>
<PublicKey 2> 
... 
<PublicKey N> 
OP_N 
OP_CHECKMULTISIG


解锁脚本

OP_0            # OP_CHECKMULTISIG的一个bug，下面会解释
<Sig1>
<Sig2>
...
<SigM>          # N个人中随意M个人的签名

完整脚本

OP_0            
<Sig1>
<Sig2>
...
<SigM>
----------
OP_M
<PublicKey 1>
<PublicKey 2> 
... 
<PublicKey N> 
OP_N 
OP_CHECKMULTISIG

```

### 脚本哈希支付P2SH（Pay to Script Hash）

```bash

锁定脚本（locking Script）形式是固定的

OP_HASH160
<redeem Script hash>   #redeem Script的哈希
OP_EQUAL

解锁脚本（unlocking Script）与赎回脚本相关

···
<Sig>
···       #需要的签名以及其他内容
<serialized redeem Script> #序列化的赎回脚本，作为数据而不作为脚本语言压栈

脚本的执行过程

首先就是解锁的脚本压栈
此时栈顶就是<serialized redeem Script>
然后锁定脚本依次压栈
先得到赎回栈顶的哈希值
再与锁定脚本中的哈希值对比两者不匹配验证就失败
成功则序列化脚本会被反序列化再与栈内的剩余内容（也就是解锁脚本中剩下的所有<Sig>或者其他内容）构成完整的脚本

```

### P2PKH 

```bash
赎回脚本(redeem Script)

OP_DUP
OP_HASH160
<address>
OP_EQUALVERIFY
OP_CHECKSIG

锁定脚本（locking Script）

OP_HASH160
<redeem Script hash>   #redeem Script的哈希
OP_EQUAL

解锁脚本(unlocking Script)

<Sig>
<Pubkey>
<serialized redeem Script>


完整脚本

<Sig>
<Pubkey>
<serialized redeem Script>
----------
OP_HASH160
<redeem Script hash>   
OP_EQUAL


```

### P2MS

```bash

赎回脚本(redeem Script)

OP_M
<PublicKey 1>
<PublicKey 2> 
... 
<PublicKey N> 
OP_N 
OP_CHECKMULTISIG

锁定脚本（locking Script）

OP_HASH160
<redeem Script hash>   #redeem Script的哈希
OP_EQUAL


解锁脚本(unlocking Script)

OP_0            
<Sig1>
<Sig2>
...
<SigM>          # N个人中随意M个人的签名
<serialized redeem Script>

完整脚本

OP_0            
<Sig1>
<Sig2>
...
<SigM>          
<serialized redeem Script>
----------
OP_HASH160
<redeem Script hash>  
OP_EQUAL

```

## note

* publicKey