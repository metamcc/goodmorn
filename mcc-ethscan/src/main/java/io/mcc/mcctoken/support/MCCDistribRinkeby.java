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
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.EthFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class MCCDistribRinkeby extends MCCDistrib {
    private static final String BINARY = "0x608060405234801561001057600080fd5b50604051610c4b380380610c4b83398101806040528101908080518201929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061008b81610091640100000000026401000000009004565b5061032f565b600080600091505b600180549050821015610144576000600260006001858154811015156100bb57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508180600101925050610099565b600090505b82518110156101ce57600160026000858481518110151561016657fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050610149565b82600190805190602001906101e4929190610262565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb836040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561024a57808201518184015260208101905061022f565b505050509050019250505060405180910390a1505050565b8280548282559060005260206000209081019282156102db579160200282015b828111156102da5782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190610282565b5b5090506102e891906102ec565b5090565b61032c91905b8082111561032857600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016102f2565b5090565b90565b61090d8061033e6000396000f300608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063025e7c271461007d578063481c6a75146100ea57806398d2fb5514610141578063a0e67e2b1461020a578063eb6b192f14610276578063fa4d3698146102d1575b600080fd5b34801561008957600080fd5b506100a860048036038101908080359060200190929190505050610337565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100f657600080fd5b506100ff610375565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561014d57600080fd5b50610208600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929050505061039a565b005b34801561021657600080fd5b5061021f61052e565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610262578082015181840152602081019050610247565b505050509050019250505060405180910390f35b34801561028257600080fd5b506102b7600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105bc565b604051808215151515815260200191505060405180910390f35b3480156102dd57600080fd5b50610335600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506105dc565b005b60018181548110151561034657fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060011515600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151415156103fc57600080fd5b849150600090505b8351811015610527578173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb858381518110151561043757fe5b90602001906020020151858481518110151561044f57fe5b906020019060200201516040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156104de57600080fd5b505af11580156104f2573d6000803e3d6000fd5b505050506040513d602081101561050857600080fd5b8101908080519060200190929190505050508080600101915050610404565b5050505050565b606060018054806020026020016040519081016040528092919081815260200182805480156105b257602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610568575b5050505050905090565b60026020528060005260406000206000915054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561063757600080fd5b61064081610643565b50565b600080600091505b6001805490508210156106f65760006002600060018581548110151561066d57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550818060010192505061064b565b600090505b825181101561078057600160026000858481518110151561071857fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555080806001019150506106fb565b8260019080519060200190610796929190610814565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb836040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156107fc5780820151818401526020810190506107e1565b505050509050019250505060405180910390a1505050565b82805482825590600052602060002090810192821561088d579160200282015b8281111561088c5782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190610834565b5b50905061089a919061089e565b5090565b6108de91905b808211156108da57600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016108a4565b5090565b905600a165627a7a72305820ce3e0aaa5db30f0d96c9b5abff5fce9fe569471656550fdeb2132ea58a5ce6850029";

    protected static final HashMap<String, String> _addresses;

    static {
        _addresses = new HashMap<String, String>();
    }

    @Deprecated
    protected MCCDistribRinkeby(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected MCCDistribRinkeby(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected MCCDistribRinkeby(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected MCCDistribRinkeby(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
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

    public static RemoteCall<MCCDistrib> deploy(Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCDistrib.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<MCCDistrib> deploy(Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCDistrib.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCDistrib> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCDistrib.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCDistrib> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, List<String> _owners) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))));
        return deployRemoteCall(MCCDistrib.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
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

    @Deprecated
    public static MCCDistribRinkeby load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCDistribRinkeby(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static MCCDistribRinkeby load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCDistribRinkeby(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static MCCDistribRinkeby load(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        return new MCCDistribRinkeby(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static MCCDistribRinkeby load(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return new MCCDistribRinkeby(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    protected String getStaticDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

    public static String getPreviouslyDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

}
