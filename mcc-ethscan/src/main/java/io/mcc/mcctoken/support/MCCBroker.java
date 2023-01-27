package io.mcc.mcctoken.support;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.EthFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.ContractGasProvider;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 3.6.0.
 */
public class MCCBroker extends Contract {
    private static final String BINARY = "0x60806040523480156200001157600080fd5b506040516200115138038062001151833981018060405260208110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b828101905060208101848111156200006757600080fd5b81518560208202830111640100000000821117156200008557600080fd5b5050929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000e781620000ee640100000000026401000000009004565b506200039d565b60008090505b600180549050811015620001a2576000600260006001848154811015156200011857fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050620000f4565b5060008090505b815181101562000231576001600260008484815181101515620001c857fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050620001a9565b5080600190805190602001906200024a929190620002c8565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb816040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015620002b257808201518184015260208101905062000295565b505050509050019250505060405180910390a150565b82805482825590600052602060002090810192821562000344579160200282015b82811115620003435782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620002e9565b5b50905062000353919062000357565b5090565b6200039a91905b808211156200039657600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016200035e565b5090565b90565b610da480620003ad6000396000f3fe608060405260043610610083576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063025e7c271461008857806343e5a30914610103578063481c6a751461029c57806398d2fb55146102f3578063a0e67e2b1461046c578063eb6b192f146104d8578063fa4d369814610541575b600080fd5b34801561009457600080fd5b506100c1600480360360208110156100ab57600080fd5b8101908080359060200190929190505050610606565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561010f57600080fd5b5061029a6004803603608081101561012657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561018357600080fd5b82018360208201111561019557600080fd5b803590602001918460208302840111640100000000831117156101b757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561021757600080fd5b82018360208201111561022957600080fd5b8035906020019184602083028401116401000000008311171561024b57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610644565b005b3480156102a857600080fd5b506102b161080d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102ff57600080fd5b5061046a6004803603606081101561031657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561035357600080fd5b82018360208201111561036557600080fd5b8035906020019184602083028401116401000000008311171561038757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803590602001906401000000008111156103e757600080fd5b8201836020820111156103f957600080fd5b8035906020019184602083028401116401000000008311171561041b57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610832565b005b34801561047857600080fd5b506104816109c6565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104c45780820151818401526020810190506104a9565b505050509050019250505060405180910390f35b3480156104e457600080fd5b50610527600480360360208110156104fb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a54565b604051808215151515815260200191505060405180910390f35b34801561054d57600080fd5b506106046004803603602081101561056457600080fd5b810190808035906020019064010000000081111561058157600080fd5b82018360208201111561059357600080fd5b803590602001918460208302840111640100000000831117156105b557600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610a74565b005b60018181548110151561061557fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60011515600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151415156106a357600080fd5b600084905060008090505b8351811015610805578173ffffffffffffffffffffffffffffffffffffffff166323b872dd8686848151811015156106e257fe5b9060200190602002015186858151811015156106fa57fe5b906020019060200201516040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b1580156107bc57600080fd5b505af11580156107d0573d6000803e3d6000fd5b505050506040513d60208110156107e657600080fd5b81019080805190602001909291905050505080806001019150506106ae565b505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60011515600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514151561089157600080fd5b600083905060008090505b83518110156109bf578173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85838151811015156108cf57fe5b9060200190602002015185848151811015156108e757fe5b906020019060200201516040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561097657600080fd5b505af115801561098a573d6000803e3d6000fd5b505050506040513d60208110156109a057600080fd5b810190808051906020019092919050505050808060010191505061089c565b5050505050565b60606001805480602002602001604051908101604052809291908181526020018280548015610a4a57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610a00575b5050505050905090565b60026020528060005260406000206000915054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610acf57600080fd5b610ad881610adb565b50565b60008090505b600180549050811015610b8c57600060026000600184815481101515610b0357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050610ae1565b5060008090505b8151811015610c18576001600260008484815181101515610bb057fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050610b93565b508060019080519060200190610c2f929190610cab565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb816040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610c95578082015181840152602081019050610c7a565b505050509050019250505060405180910390a150565b828054828255906000526020600020908101928215610d24579160200282015b82811115610d235782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190610ccb565b5b509050610d319190610d35565b5090565b610d7591905b80821115610d7157600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905550600101610d3b565b5090565b9056fea165627a7a723058201ef45811cb19f25b3752cc46d2652cd6d3962fb97baa2a4e5437b76b08f246b60029";

    public static final String FUNC_OWNERS = "owners";

    public static final String FUNC_MANAGER = "manager";

    public static final String FUNC_GETOWNERS = "getOwners";

    public static final String FUNC_OWNERBYADDRESS = "ownerByAddress";

    public static final String FUNC_SETOWNERS = "setOwners";

    public static final String FUNC_TRANSFERMULTI = "transferMulti";

    public static final String FUNC_TRANSFERFROMMULTI = "transferFromMulti";

    public static final Event SETOWNERS_EVENT = new Event("SetOwners", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
    ;

    protected static final HashMap<String, String> _addresses;

    static {
        _addresses = new HashMap<String, String>();
        _addresses.put("1", "0x3bA9174C9a4251bb33b43619eC7c7519779C52dE");
        _addresses.put("4", "0x2DEcb775466bC8eC3508b292579a10B274bB9CBE");
    }

    @Deprecated
    protected MCCBroker(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected MCCBroker(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected MCCBroker(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected MCCBroker(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> owners(BigInteger param0) {
        final Function function = new Function(FUNC_OWNERS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> manager() {
        final Function function = new Function(FUNC_MANAGER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<List> getOwners() {
        final Function function = new Function(FUNC_GETOWNERS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<Boolean> ownerByAddress(String param0) {
        final Function function = new Function(FUNC_OWNERBYADDRESS, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> setOwners(List<String> _owners) {
        final Function function = new Function(
                FUNC_SETOWNERS, 
                Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MCCBroker> deploy(Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCBroker.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<MCCBroker> deploy(Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCBroker.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCBroker> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCBroker.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCBroker> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCBroker.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    public List<SetOwnersEventResponse> getSetOwnersEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(SETOWNERS_EVENT, transactionReceipt);
        ArrayList<SetOwnersEventResponse> responses = new ArrayList<SetOwnersEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            SetOwnersEventResponse typedResponse = new SetOwnersEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.owners = (List<String>) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<SetOwnersEventResponse> setOwnersEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, SetOwnersEventResponse>() {
            @Override
            public SetOwnersEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(SETOWNERS_EVENT, log);
                SetOwnersEventResponse typedResponse = new SetOwnersEventResponse();
                typedResponse.log = log;
                typedResponse.owners = (List<String>) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<SetOwnersEventResponse> setOwnersEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(SETOWNERS_EVENT));
        return setOwnersEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> transferMulti(String _tokenAddr, List<String> targets, List<BigInteger> values) {
        final Function function = new Function(
                FUNC_TRANSFERMULTI, 
                Arrays.<Type>asList(new Address(_tokenAddr),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(targets, Address.class)),
                new DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(values, org.web3j.abi.datatypes.generated.Uint256.class))), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> transferFromMulti(String _tokenAddr, String fromAddr, List<String> targets, List<BigInteger> values) {
        final Function function = new Function(
                FUNC_TRANSFERFROMMULTI, 
                Arrays.<Type>asList(new Address(_tokenAddr),
                new Address(fromAddr),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(targets, Address.class)),
                new DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(values, org.web3j.abi.datatypes.generated.Uint256.class))), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    @Deprecated
    public static MCCBroker load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCBroker(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static MCCBroker load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCBroker(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static MCCBroker load(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        return new MCCBroker(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static MCCBroker load(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return new MCCBroker(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    protected String getStaticDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

    public static String getPreviouslyDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

    public static class SetOwnersEventResponse {
        public Log log;

        public List<String> owners;
    }
}
