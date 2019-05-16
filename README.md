# Project 5 

## A Blockchain Based Voting System
### Why: 
   There is human error as well as corruption in the voting process, and vulnerabilities exist throughout the current voting process from start to end. The transportation and counting of paper ballots leaves the process of counting and recording votes open to serious vulnerabilities. Paper ballots can be directly modified, and whole boxes of votes have been lost or stolen in certain districts. Also, the digitized voting process would be cheaper long term. Cost of paper, mail, and labor for managing polling stations in the current infrastructure will be more than the creation and maintenance of this system. There are similar issues with states that have implemented voting machines, and poor implementations of such machines have already decimated public trust in them.
   
   Digital voting is not feasible without blockchain technology. Having a publicly readable, immutable and verifiable ledger will provide transparency, and thus interest and faith in the voting process. Having a decentralized network of nodes (every polling station can be replaced by a node) confirms that the data cannot be tampered with. In a centralized system, this is not the case. If the one centralized system is compromised, all the voting results are also compromised. 
### How: 
To effectively digitalize the voting process, the system must meet these criteria:

* Integrity: Only eligible voters may vote, and they may only vote once. (Thus at some point there will need to be verification off the chain).
The blockchain will need to have a permissioned ledger so only permissioned nodes (polling stations) have write permissions. Read permissions will be public.
* Transparency: All results should be independently verifiable.
* Privacy: Choices of a voter must be kept private during and after the election. (Prevent voting under duress/buying of votes).

The voter registration process would still have to take place off the grid or off the chain. Once a voter registration agency determines that someone is eligible to vote, they would receive a token or key that would allow them to vote exactly once. This can be distributed digitally, or via snail mail. Either way, since your identity needs to be confirmed at some point during the entire voting process, this is unavoidable. This guarantees the integrity requirement. Also, it is assumed that polling stations will still exist, and voters would have to go to the stations in order to vote. They would just do so digitally. Although it would be ideal to allow the user to vote from the comfort of their own home, that is outside the scope of this project.
### Project Description
Develop the voting system detailed above with the fundamental requirements of Transparency and Privacy. Integrity will be partially guaranteed off chain, and partially guaranteed by having a ledger with write permissions.  

Send: A node will await user input of an ID string from a node allowed (determined by the permissioned ledger). If that ID is confirmed as correct (make an API call to a centralized server that verifies user identity) and the node is permissioned, the node will connect to the network and download the voting ballot. The node will provide a prompt for inputting votes. User will make their decisions and save/confirm. The voting results will be hashed, and the block will be encrypted with the ring signature to maintain anonymity. That will be pushed to the blockchain. 

Receive: The other nodes will need to verify the new block using the key image.
Success Criteria:
* Only specific nodes/roles can write to the blockchain (Permissioned Ledger successfully implemented)
* Blockchain does not fork, or forks are resolved immediately and no data is lost.
* Voting nodes can submit anonymous votes to the blockchain (Ring Signature successfully implemented)
* Votes can be verified on the blockchain (but identities of the submitter remain anonymous)

## Final Timeline

Link to Demo: https://www.youtube.com/watch?v=jlRtCUJsp9k

Week 1a: Grant write privileges to blockchain from only certain nodes (Permissioned Ledger) :white_check_mark:

Week 1b: Provide voting user interface and downloadable ballet :white_check_mark:

Week 2a: Implement signed data transactions, public and private keys :white_check_mark:

Week 2b: Validate digital signature :white_check_mark:

Week 4a: Implement solution to blockchain forking. Count votes on all chains so forks don't impact final count. :white_check_mark:

Week 4b: Verify integrity of the votes and count the votes :white_check_mark:

Week 4b: Read in the ballot from JSON file :white_check_mark:

Week 4c: Conceal signature/public key of voter with ring signature algorithm :hourglass_flowing_sand:

Week 4c: Validate ring signature  :hourglass_flowing_sand:

Final Due Date/Demo - 5/16/19 at 2:30 p.m.

### Final Review
Did not successfully implement ring signature by deadline. The bug for verifying the signature took longer than expected to resolve 
(ended up using the go-ethereum library instead of the ECDSA library, since it wasn't clear how the receiving node would validate the signature.) 



## Midpoint Review

The original timeline did not allocate time to implement public/private keys or digital signatures.
Ring signature functionality is dependent on having some type of existing digital signature infrastructure in place.
I've moved back "Implement solution to blockchain forking" but was still able to start on the ring signature on schedule.


### Updated Timeline:

Week 1a: Grant write privileges to blockchain from only certain nodes (Permissioned Ledger) :white_check_mark:

Week 1b: Provide voting user interface and downloadable ballet :white_check_mark:

Week 2a: Implement signed data transactions, public and private keys :white_check_mark:

Week 2b: Verify on the blockchain by digital signature :hourglass_flowing_sand: (Currently a bug in verification)

CHECKPOINT - 5/1/19

Week 3a: Start work on encrypting voting results using ring signature algorithm :hourglass_flowing_sand:

Week 3b: Validate ring signature results  :x:

Week 4a: Implement solution to blockchain forking. (Implement Round Robin Consensus, Delegated Proof of Stake, etc. or use a Canononical Tree)  :x:

Week 4b: Verify integrity of the votes and count the votes  :x:
Final Due Date/Demo - 5/16/19 at 2:30 p.m.

### Permissioned Blockchain
<addr>
    
    /* For simplicity's sake, this can just be the Port Numbers since we're collecting that info already.
    In production, we can actually keep a seperate list of predetermined allowed (Public?) IDs */
    
     /* Only these IDs are allowed to write. This gives the semblance of a permissioned blockchain */
    var ALLOWED_IDS = map[int32]bool {
        6688:true,
        6669:true,
        6670:true,
    }
    
</addr>

### User Interface for Voting
Added a CLI so that the voter can submit their vote. The block creation blocks while the user votes.
<addr>

    func GenerateVotingMpt() p1.MerklePatriciaTrie {
        mpt := p1.MerklePatriciaTrie{}
        mpt.Initial()
        scanner := bufio.NewScanner(os.Stdin)
        var text string
        for text != "q" {  // break the loop if text == "q"
            fmt.Print("Enter your Vote: ")
            scanner.Scan()
            text = scanner.Text()
            if text != "q" {
                fmt.Println("You voted for ", text)
            }
        }
        fmt.Println("Thanks for voting!")
        /* Record the Vote in the MPT. Value can be vote value. Key will be the
        public key i think.
        */
        mpt.Insert(text, text)
    
        if scanner.Err() != nil {
            // handle error.
        }
        return mpt
    }

</addr>

### Downloadable Ballot
The node must download the ballot from the blockchain so they know their voting choices.
<addr>

    /* POST the contents of the ballot so other nodes can download. Read from ballot.json */
    func UploadBallot(w http.ResponseWriter, r *http.Request) {
        /* Read the JSON File */
        plan, _ := ioutil.ReadFile("ballot.json")
        var data interface{}
        err := json.Unmarshal(plan, &data)
        if err != nil {
            fmt.Println("Error")
        }
        fmt.Fprint(w, plan)
        fmt.Println(plan)
    }
    
    
</addr>

### Data Transactions and Digital Signature
To implement a ring signature and guarantee integrity of data, the first step is to implement digital signatures.
The current implementation of SignTransation() takes the MPT root as input, and signs
that value with the node's private key.

The current implementation of VerifySignature() takes as input a block that's ready
to be added to the blockchain, and outputs a boolean.
<addr>

    func SignTransaction(value string) {
        transaction := []byte (value)
        r := big.NewInt(0)
        s := big.NewInt(0)
        serr := errors.New("Error")
        /* Returns Big Ints r and s
        */
        r, s, serr = ecdsa.Sign(crand.Reader, PRIVATE_KEY, transaction)
        if serr != nil {
            fmt.Println("Error")
            os.Exit(1)
        }
    
        /* Need to figure out how to get the r and s values from the signature on the blockchain */
    
        /*
        The signature is a combination of the author's private key and the content of
        the document it certifies
         */
        SIGNATURE = r.Bytes()
        SIGNATURE = append(SIGNATURE, s.Bytes()...)
    }
      
    func VerifySignature(block p2.Block) bool{
    	e := &ECDSASignature{}
    	_, err := asn1.Unmarshal([]byte(block.Header.Signature), e)
    	if err != nil {
    		fmt.Println("Error Unmarshaling Block")
    		return false
    	}
    	verified := ecdsa.Verify(block.PublicKey, []byte(block.GetMptRoot()), e.R, e.S)
    	return verified
    }
    
</addr>

### Public and Private Key Generation
We use the ECDSA golang library to generate our Public and Private keys.

<addr>

    func GeneratePublicAndPrivateKey() {
        c := elliptic.P256()
        PRIVATE_KEY, _ = ecdsa.GenerateKey(c, crand.Reader)
        PUBLIC_KEY = append(PRIVATE_KEY.PublicKey.X.Bytes(), PRIVATE_KEY.PublicKey.Y.Bytes()...)
    }
    
    /* Hash the public key to display on the blockchain. This is how BTC does it */
    func HashPublicKey(publicKey []byte) []byte {
        publicSHA256 := sha256.Sum256(publicKey)
        ripemd160Hasher := ripemd160.New()
        _, err := ripemd160Hasher.Write(publicSHA256[:])
    
        if err != nil {
            fmt.Println("Cannot Hash Public Key")
            os.Exit(1)
        }
    
        hashedPublicKey := ripemd160Hasher.Sum(nil)
        return hashedPublicKey
    }

</addr>

## Original Timeline

Week 1a: Grant write privileges to blockchain from only certain nodes (Permissioned Ledger)

Week 1b: Provide voting user interface and downloadable ballet

Week 2a: Implement solution to blockchain forking. (Implement Round Robin Consensus, Delegated Proof of Stake, or something similar)

Week 2b: Start work on encrypting voting results using ring signature algorithm

CHECKPOINT - 5/1/19

Week 3a: Complete ring signature algorithm

Week 3b: Verify on the blockchain by the key image

Week 4: Verify integrity of the votes and count the votes


## Resources

### Permissioned Blockchain
* https://www.coindesk.com/information/what-is-the-difference-between-open-and-permissioned-blockchains
* https://medium.com/coinmonks/permissioned-blockchains-are-a-dead-end-67c2b060bc52
* https://monax.io/learn/permissioned_blockchains/

### Transactions
* https://jeiwan.cc/posts/building-blockchain-in-go-part-4/

### Digital Signature
* https://golang.org/pkg/crypto/ecdsa
* https://lisk.io/academy/blockchain-basics/how-does-blockchain-work/digital-signatures
* https://medium.com/@xragrawal/digital-signature-from-blockchain-context-cedcd563eee5
* https://medium.com/icovo/digital-signatures-in-a-blockchain-digital-signatures-44b981b75413

### Ring Signature
* https://www.mycryptopedia.com/monero-ring-signature-explained/
* https://en.wikipedia.org/wiki/Ring_signature
* https://blockonomi.com/ring-signatures/

### Different Methods of Consensus
* https://blockchain.intellectsoft.net/blog/consensus-protocols-that-serve-different-business-needs-part-2/
* https://blockchainlion.com/consensus-blockchain/

### Voting
* https://static1.squarespace.com/static/5b0be2f4e2ccd12e7e8a9be9/t/5b6c38550e2e725e9cad3f18/1533818968655/Agora_Whitepaper.pdf

## Documentation
You will need to launch the initial node on PORT 6688, and then send a GET request to /create the create the first blockchain. 
You will then need to send a GET request to /start to start creating blocks. You can then run and start the rest of the nodes.