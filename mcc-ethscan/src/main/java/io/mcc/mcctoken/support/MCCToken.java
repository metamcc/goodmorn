package io.mcc.mcctoken.support;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.EthFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.ContractGasProvider;
import org.web3j.tx.gas.StaticGasProvider;
import rx.Observable;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 3.6.0.
 */
public abstract class MCCToken extends Contract {

    private static final String BINARY = null;

    public static final String FUNC_OWNERS = "owners";

    public static final String FUNC_NAME = "name";

    public static final String FUNC_APPROVE = "approve";

    public static final String FUNC_TOTALSUPPLY = "totalSupply";

    public static final String FUNC_ALLOWTRANSFERS = "allowTransfers";

    public static final String FUNC_TRANSFERFROM = "transferFrom";

    public static final String FUNC_LIMITEDWALLETS = "limitedWallets";

    public static final String FUNC_DECIMALS = "decimals";

    public static final String FUNC_LIMITEDWALLETSMANAGER = "limitedWalletsManager";

    public static final String FUNC_ISSUANCEFINISHED = "issuanceFinished";

    public static final String FUNC_MANAGER = "manager";

    public static final String FUNC_DECREASEAPPROVAL = "decreaseApproval";

    public static final String FUNC_BALANCEOF = "balanceOf";

    public static final String FUNC_DELLIMITEDWALLETADDRESS = "delLimitedWalletAddress";

    public static final String FUNC_SYMBOL = "symbol";

    public static final String FUNC_SETLIMITENABLED = "setLimitEnabled";

    public static final String FUNC_GETOWNERS = "getOwners";

    public static final String FUNC_DESTROY = "destroy";

    public static final String FUNC_ISLIMITEDWALLETADDRESS = "isLimitedWalletAddress";

    public static final String FUNC_TRANSFER = "transfer";

    public static final String FUNC_SETLISTENER = "setListener";

    public static final String FUNC_FINISHISSUANCE = "finishIssuance";

    public static final String FUNC_EVENTLISTENER = "eventListener";

    public static final String FUNC_INCREASEAPPROVAL = "increaseApproval";

    public static final String FUNC_ISLIMITENABLED = "isLimitEnabled";

    public static final String FUNC_ALLOWANCE = "allowance";

    public static final String FUNC_SETALLOWTRANSFERS = "setAllowTransfers";

    public static final String FUNC_OWNERBYADDRESS = "ownerByAddress";

    public static final String FUNC_ADDLIMITEDWALLETADDRESS = "addLimitedWalletAddress";

    public static final String FUNC_SETOWNERS = "setOwners";

    public static final String FUNC_ISSUE = "issue";

    public static final String FUNC_BURN = "burn";

    public static final String FUNC_BURNFROM = "burnFrom";

    public static final Event BURN_EVENT = new Event("Burn",
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event ALLOWTRANSFERSCHANGED_EVENT = new Event("AllowTransfersChanged",
            Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
    ;

    public static final Event ISSUE_EVENT = new Event("Issue",
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event DESTROY_EVENT = new Event("Destroy",
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event ISSUANCEFINISHED_EVENT = new Event("IssuanceFinished",
            Arrays.<TypeReference<?>>asList());
    ;

    public static final Event SETOWNERS_EVENT = new Event("SetOwners",
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
    ;

    public static final Event TRANSFER_EVENT = new Event("Transfer",
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event APPROVAL_EVENT = new Event("Approval",
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>(true) {}, new TypeReference<Address>(true) {}, new TypeReference<Uint256>() {}));
    ;

    protected MCCToken(String contractBinary, String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider gasProvider) {
        super(contractBinary, contractAddress, web3j, transactionManager, gasProvider);
    }

    protected MCCToken(String contractBinary, String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider gasProvider) {
        super(contractBinary, (String)contractAddress, (Web3j)web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), (ContractGasProvider)gasProvider);
    }

    protected MCCToken(String contractBinary, String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractBinary, (String)contractAddress, (Web3j)web3j, (TransactionManager)transactionManager, (ContractGasProvider)(new StaticGasProvider(gasPrice, gasLimit)));
    }

    protected MCCToken(String contractBinary, String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractBinary, contractAddress, web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), gasPrice, gasLimit);
    }

    protected MCCToken(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super("", contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

   protected MCCToken(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
       super("", contractAddress, web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), gasPrice, gasLimit);
    }


    public abstract RemoteCall<String> owners(BigInteger param0);

    public abstract RemoteCall<String> name();

    public abstract RemoteCall<TransactionReceipt> approve(String _spender, BigInteger _value);

    public abstract RemoteCall<BigInteger> totalSupply();
    public abstract RemoteCall<Boolean> allowTransfers();

    public abstract RemoteCall<TransactionReceipt> transferFrom(String _from, String _to, BigInteger _value);

    public abstract RemoteCall<Boolean> limitedWallets(String param0) ;

    public abstract RemoteCall<BigInteger> decimals();

    public abstract RemoteCall<String> limitedWalletsManager() ;

    public abstract RemoteCall<Boolean> issuanceFinished();
    public abstract RemoteCall<String> manager() ;

    public abstract RemoteCall<TransactionReceipt> decreaseApproval(String _spender, BigInteger _subtractedValue);

    public abstract RemoteCall<BigInteger> balanceOf(String _owner);

    public abstract RemoteCall<TransactionReceipt> delLimitedWalletAddress(String _wallet);
    public abstract RemoteCall<String> symbol() ;

    public abstract RemoteCall<TransactionReceipt> setLimitEnabled(Boolean _setLimitEnabled);

    public abstract RemoteCall<List> getOwners();

    public abstract RemoteCall<TransactionReceipt> destroy(String _from, BigInteger _value);

    public abstract RemoteCall<Boolean> isLimitedWalletAddress(String _wallet);

    public abstract RemoteCall<TransactionReceipt> transfer(String _to, BigInteger _value);

    public abstract RemoteCall<TransactionReceipt> setListener(String _listener);

    public abstract RemoteCall<TransactionReceipt> finishIssuance();

    public abstract RemoteCall<String> eventListener();
    public abstract RemoteCall<TransactionReceipt> increaseApproval(String _spender, BigInteger _addedValue);

    public abstract RemoteCall<Boolean> isLimitEnabled();

    public abstract RemoteCall<BigInteger> allowance(String _owner, String _spender);

    public abstract RemoteCall<TransactionReceipt> setAllowTransfers(Boolean _allowTransfers);

    public abstract RemoteCall<Boolean> ownerByAddress(String param0);

    public abstract RemoteCall<TransactionReceipt> addLimitedWalletAddress(String _wallet);

    public abstract RemoteCall<TransactionReceipt> setOwners(List<String> _owners);

    public abstract List<BurnEventResponse> getBurnEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<BurnEventResponse> burnEventObservable(EthFilter filter);

    public abstract Observable<BurnEventResponse> burnEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

    public abstract List<AllowTransfersChangedEventResponse> getAllowTransfersChangedEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<AllowTransfersChangedEventResponse> allowTransfersChangedEventObservable(EthFilter filter);

    public abstract Observable<AllowTransfersChangedEventResponse> allowTransfersChangedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) ;

    public abstract List<IssueEventResponse> getIssueEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<IssueEventResponse> issueEventObservable(EthFilter filter);

    public abstract Observable<IssueEventResponse> issueEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

    public abstract List<DestroyEventResponse> getDestroyEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<DestroyEventResponse> destroyEventObservable(EthFilter filter);

    public abstract Observable<DestroyEventResponse> destroyEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

    public abstract List<IssuanceFinishedEventResponse> getIssuanceFinishedEvents(TransactionReceipt transactionReceipt);
    public abstract Observable<IssuanceFinishedEventResponse> issuanceFinishedEventObservable(EthFilter filter);

    public abstract Observable<IssuanceFinishedEventResponse> issuanceFinishedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

    public abstract List<SetOwnersEventResponse> getSetOwnersEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<SetOwnersEventResponse> setOwnersEventObservable(EthFilter filter);

    public abstract Observable<SetOwnersEventResponse> setOwnersEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);
    public abstract List<TransferEventResponse> getTransferEvents(TransactionReceipt transactionReceipt);
    public abstract Observable<TransferEventResponse> transferEventObservable(EthFilter filter);

    public abstract Observable<TransferEventResponse> transferEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);
    public abstract List<ApprovalEventResponse> getApprovalEvents(TransactionReceipt transactionReceipt);

    public abstract Observable<ApprovalEventResponse> approvalEventObservable(EthFilter filter);

    public abstract Observable<ApprovalEventResponse> approvalEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

    public abstract RemoteCall<TransactionReceipt> issue(String _to, BigInteger _value);

    public abstract RemoteCall<TransactionReceipt> burn(BigInteger _value);

    public abstract RemoteCall<TransactionReceipt> burnFrom(String _from, BigInteger _value);

    public static class BurnEventResponse {
        public Log log;

        public String burner;

        public BigInteger value;
    }

    public static class AllowTransfersChangedEventResponse {
        public Log log;

        public Boolean _newState;
    }

    public static class IssueEventResponse {
        public Log log;

        public String _to;

        public BigInteger _value;
    }

    public static class DestroyEventResponse {
        public Log log;

        public String _from;

        public BigInteger _value;
    }

    public static class IssuanceFinishedEventResponse {
        public Log log;
    }

    public static class SetOwnersEventResponse {
        public Log log;

        public List<String> owners;
    }

    public static class TransferEventResponse {
        public Log log;

        public String from;

        public String to;

        public BigInteger value;
    }

    public static class ApprovalEventResponse {
        public Log log;

        public String owner;

        public String spender;

        public BigInteger value;
    }
}
