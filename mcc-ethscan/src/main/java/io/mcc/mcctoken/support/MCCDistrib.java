package io.mcc.mcctoken.support;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
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
public abstract class MCCDistrib extends Contract {

    private static final String BINARY = null;

    public static final String FUNC_OWNERS = "owners";

    public static final String FUNC_MANAGER = "manager";

    public static final String FUNC_GETOWNERS = "getOwners";

    public static final String FUNC_OWNERBYADDRESS = "ownerByAddress";

    public static final String FUNC_SETOWNERS = "setOwners";

    public static final String FUNC_TRANSFERMULTI = "transferMulti";

    public static final Event SETOWNERS_EVENT = new Event("SetOwners",
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
    ;

    protected MCCDistrib(String contractBinary, String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider gasProvider) {
        super(contractBinary, contractAddress, web3j, transactionManager, gasProvider);
    }

    protected MCCDistrib(String contractBinary, String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider gasProvider) {
        super(contractBinary, (String)contractAddress, (Web3j)web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), (ContractGasProvider)gasProvider);
    }

    protected MCCDistrib(String contractBinary, String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractBinary, (String)contractAddress, (Web3j)web3j, (TransactionManager)transactionManager, (ContractGasProvider)(new StaticGasProvider(gasPrice, gasLimit)));
    }

    protected MCCDistrib(String contractBinary, String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractBinary, contractAddress, web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), gasPrice, gasLimit);
    }

    protected MCCDistrib(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super("", contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

   protected MCCDistrib(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
       super("", contractAddress, web3j, (TransactionManager)(new RawTransactionManager(web3j, credentials)), gasPrice, gasLimit);
    }


   public abstract RemoteCall<String> owners(BigInteger param0);

   public abstract RemoteCall<String> manager();

   public abstract RemoteCall<List> getOwners();

   public abstract RemoteCall<Boolean> ownerByAddress(String param0);

   public abstract RemoteCall<TransactionReceipt> setOwners(List<String> _owners);

   public abstract List<SetOwnersEventResponse> getSetOwnersEvents(TransactionReceipt transactionReceipt);

   public abstract Observable<SetOwnersEventResponse> setOwnersEventObservable(EthFilter filter);

   public abstract Observable<SetOwnersEventResponse> setOwnersEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock);

   public abstract RemoteCall<TransactionReceipt> transferMulti(String _tokenAddr, List<String> targets, List<BigInteger> values);

   protected abstract String getStaticDeployedAddress(String networkId);

    public static class SetOwnersEventResponse {
        public Log log;

        public List<String> owners;
    }
}
