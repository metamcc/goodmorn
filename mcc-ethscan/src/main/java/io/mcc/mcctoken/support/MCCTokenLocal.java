package io.mcc.mcctoken.support;

import org.web3j.abi.EventEncoder;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.*;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
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

import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 3.6.0.
 */
public class MCCTokenLocal extends MCCToken {
    private static final String BINARY = "0x60806040526000600660006101000a81548160ff0219169083151502179055506000600660016101000a81548160ff0219169083151502179055506040805190810160405280600d81526020017f4d79437265646974436861696e00000000000000000000000000000000000000815250600990805190602001906200008792919062000535565b506040805190810160405280600381526020017f4d43430000000000000000000000000000000000000000000000000000000000815250600a9080519060200190620000d592919062000535565b506012600b60006101000a81548160ff021916908360ff160217905550348015620000ff57600080fd5b50604051620035d7380380620035d7833981018060405281019080805190602001909291908051820192919060200180519060200190929190505050828282828233600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141515620001fa5781600660026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b62000214816200031e640100000000026401000000009004565b50506001600860146101000a81548160ff02191690831515021790555080600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050620002b1600b60009054906101000a900460ff1660ff16600a0a633b9aca00620004f96401000000000262002eaa179091906401000000009004565b600281905550600254600080846000815181101515620002cd57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505050620006b9565b600080600091505b600480549050821015620003d4576000600560006004858154811015156200034a57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550818060010192505062000326565b600090505b825181101562000461576001600560008584815181101515620003f857fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050620003d9565b826004908051906020019062000479929190620005bc565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb836040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015620004e1578082015181840152602081019050620004c4565b505050509050019250505060405180910390a1505050565b6000808314156200050e57600090506200052f565b81830290508183828115156200052057fe5b041415156200052b57fe5b8090505b92915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200057857805160ff1916838001178555620005a9565b82800160010185558215620005a9579182015b82811115620005a85782518255916020019190600101906200058b565b5b509050620005b891906200064b565b5090565b82805482825590600052602060002090810192821562000638579160200282015b82811115620006375782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620005dd565b5b50905062000647919062000673565b5090565b6200067091905b808211156200066c57600081600090555060010162000652565b5090565b90565b620006b691905b80821115620006b257600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016200067a565b5090565b90565b612f0e80620006c96000396000f3006080604052600436106101a1576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063025e7c27146101a657806306fdde0314610213578063095ea7b3146102a357806318160ddd146103085780632185810b1461033357806323b872dd146103625780632e217405146103e7578063313ce5671461044257806342966c681461047357806344e7faa4146104a05780634662299a146104f7578063481c6a7514610526578063661884631461057d57806370a08231146105e257806379cc6790146106395780637d80265514610686578063867904b4146106c957806395d89b41146107165780639da793d0146107a6578063a0e67e2b146107d5578063a24835d114610841578063a24ed4e51461088e578063a9059cbb146108e9578063adcd905b1461094e578063c422293b14610991578063cd9217f7146109a8578063d73dd623146109ff578063daf4f66e14610a64578063dd62ed3e14610a93578063df50afa414610b0a578063eb6b192f14610b39578063ee8cbc9d14610b94578063fa4d369814610bd7575b600080fd5b3480156101b257600080fd5b506101d160048036038101908080359060200190929190505050610c3d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561021f57600080fd5b50610228610c7b565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561026857808201518184015260208101905061024d565b50505050905090810190601f1680156102955780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156102af57600080fd5b506102ee600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610d19565b604051808215151515815260200191505060405180910390f35b34801561031457600080fd5b5061031d610df8565b6040518082815260200191505060405180910390f35b34801561033f57600080fd5b50610348610e02565b604051808215151515815260200191505060405180910390f35b34801561036e57600080fd5b506103cd600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610e15565b604051808215151515815260200191505060405180910390f35b3480156103f357600080fd5b50610428600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ef6565b604051808215151515815260200191505060405180910390f35b34801561044e57600080fd5b50610457610f16565b604051808260ff1660ff16815260200191505060405180910390f35b34801561047f57600080fd5b5061049e60048036038101908080359060200190929190505050610f29565b005b3480156104ac57600080fd5b506104b5610f36565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561050357600080fd5b5061050c610f5c565b604051808215151515815260200191505060405180910390f35b34801561053257600080fd5b5061053b610f6f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561058957600080fd5b506105c8600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610f95565b604051808215151515815260200191505060405180910390f35b3480156105ee57600080fd5b50610623600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061121d565b6040518082815260200191505060405180910390f35b34801561064557600080fd5b50610684600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611265565b005b34801561069257600080fd5b506106c7600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061140d565b005b3480156106d557600080fd5b50610714600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506114c4565b005b34801561072257600080fd5b5061072b611707565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561076b578082015181840152602081019050610750565b50505050905090810190601f1680156107985780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156107b257600080fd5b506107d36004803603810190808035151590602001909291905050506117a5565b005b3480156107e157600080fd5b506107ea61181e565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561082d578082015181840152602081019050610812565b505050509050019250505060405180910390f35b34801561084d57600080fd5b5061088c600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506118ac565b005b34801561089a57600080fd5b506108cf600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506118b0565b604051808215151515815260200191505060405180910390f35b3480156108f557600080fd5b50610934600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611906565b604051808215151515815260200191505060405180910390f35b34801561095a57600080fd5b5061098f600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506119e5565b005b34801561099d57600080fd5b506109a6611ae9565b005b3480156109b457600080fd5b506109bd611b91565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015610a0b57600080fd5b50610a4a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611bb7565b604051808215151515815260200191505060405180910390f35b348015610a7057600080fd5b50610a79611daa565b604051808215151515815260200191505060405180910390f35b348015610a9f57600080fd5b50610af4600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611dbd565b6040518082815260200191505060405180910390f35b348015610b1657600080fd5b50610b37600480360381019080803515159060200190929190505050611e44565b005b348015610b4557600080fd5b50610b7a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611efb565b604051808215151515815260200191505060405180910390f35b348015610ba057600080fd5b50610bd5600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611f1b565b005b348015610be357600080fd5b50610c3b60048036038101908080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290505050611fd2565b005b600481815481101515610c4c57fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60098054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d115780601f10610ce657610100808354040283529160200191610d11565b820191906000526020600020905b815481529060010190602001808311610cf457829003601f168201915b505050505081565b60003383600860149054906101000a900460ff161580610dd95750600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16158015610dd85750600760008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16155b5b1515610de457600080fd5b610dee858561203a565b9250505092915050565b6000600254905090565b600660009054906101000a900460ff1681565b60008383600860149054906101000a900460ff161580610ed55750600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16158015610ed45750600760008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16155b5b1515610ee057600080fd5b610eeb86868661212c565b925050509392505050565b60076020528060005260406000206000915054906101000a900460ff1681565b600b60009054906101000a900460ff1681565b610f333382612287565b50565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600660019054906101000a900460ff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050808311156110a6576000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550611131565b6110b0818461243a565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505b8373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040518082815260200191505060405180910390a3600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205481111515156112f057600080fd5b61137f81600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461243a90919063ffffffff16565b600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506114098282612287565b5050565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561146957600080fd5b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b60011515600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514151561152357600080fd5b600660019054906101000a900460ff1615151561153c57fe5b6115846000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548261243a565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061160e6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482612453565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff167fc65a3f767206d2fdcede0b094a4840e01c0dd0be1888b5ba800346eaa0123c16826040518082815260200191505060405180910390a28173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b600a8054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561179d5780601f106117725761010080835404028352916020019161179d565b820191906000526020600020905b81548152906001019060200180831161178057829003601f168201915b505050505081565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561180157600080fd5b80600860146101000a81548160ff02191690831515021790555050565b606060048054806020026020016040519081016040528092919081815260200182805480156118a257602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311611858575b5050505050905090565b5050565b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b60003383600860149054906101000a900460ff1615806119c65750600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161580156119c55750600760008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16155b5b15156119d157600080fd5b6119db858561246f565b9250505092915050565b60011515600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515611a4457600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515611ac05780600660026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611ae6565b600660026101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690555b50565b60011515600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515611b4857600080fd5b6001600660016101000a81548160ff0219169083151502179055507f29fe76cc5ca143e91eadf7242fda487fcef09318c1237900f958abe1e2c5beff60405160405180910390a1565b600660029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000611c3f600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205483612453565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040518082815260200191505060405180910390a36001905092915050565b600860149054906101000a900460ff1681565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60011515600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515611ea357600080fd5b80600660006101000a81548160ff0219169083151502179055507fbac956a1816a25b65e25a2449379c8409891b96663ce5f0b3475c196ec4bfa0f81604051808215151515815260200191505060405180910390a150565b60056020528060005260406000206000915054906101000a900460ff1681565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515611f7757600080fd5b6001600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561202e57600080fd5b612037816125c8565b50565b600081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b600080600660009054906101000a900460ff16151561214757fe5b612152858585612799565b905061215c612b54565b80156121655750805b1561227c57600660029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663677ba3d38686866040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050600060405180830381600087803b15801561226357600080fd5b505af1158015612277573d6000803e3d6000fd5b505050505b809150509392505050565b6000808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205481111515156122d457600080fd5b612325816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461243a90919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061237c8160025461243a90919063ffffffff16565b6002819055508173ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5826040518082815260200191505060405180910390a2600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b600082821115151561244857fe5b818303905092915050565b6000818301905082811015151561246657fe5b80905092915050565b600080600660009054906101000a900460ff16151561248a57fe5b6124948484612bbd565b905061249e612b54565b80156124a75750805b156125be57600660029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663677ba3d33386866040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050600060405180830381600087803b1580156125a557600080fd5b505af11580156125b9573d6000803e3d6000fd5b505050505b8091505092915050565b600080600091505b60048054905082101561267b576000600560006004858154811015156125f257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555081806001019250506125d0565b600090505b825181101561270557600160056000858481518110151561269d57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050612680565b826004908051906020019061271b929190612ddd565b507f9465cd279c2de393c5568ae444599e3644e3d1864ca2c05ced8a654df2aea3cb836040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015612781578082015181840152602081019050612766565b505050509050019250505060405180910390a1505050565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482111515156127e857600080fd5b600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821115151561287357600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141515156128af57600080fd5b612900826000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461243a90919063ffffffff16565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550612993826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461245390919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550612a6482600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461243a90919063ffffffff16565b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190509392505050565b60008073ffffffffffffffffffffffffffffffffffffffff16600660029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415612bb55760009050612bba565b600190505b90565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548211151515612c0c57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515612c4857600080fd5b612c99826000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461243a90919063ffffffff16565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550612d2c826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461245390919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b828054828255906000526020600020908101928215612e56579160200282015b82811115612e555782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190612dfd565b5b509050612e639190612e67565b5090565b612ea791905b80821115612ea357600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905550600101612e6d565b5090565b90565b600080831415612ebd5760009050612edc565b8183029050818382811515612ece57fe5b04141515612ed857fe5b8090505b929150505600a165627a7a7230582004ffb623fcb6f46be04ed8154ce5de0208270e6664fb2dfaadeb106f251c23c20029";

    protected static final HashMap<String, String> _addresses;

    static {
        _addresses = new HashMap<String, String>();
        _addresses.put("4649", "0x56d9848bd2fd5b971fbe687ca57e3e0cf9357324");
    }

    @Deprecated
    protected MCCTokenLocal(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected MCCTokenLocal(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected MCCTokenLocal(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected MCCTokenLocal(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> owners(BigInteger param0) {
        final Function function = new Function(FUNC_OWNERS, 
                Arrays.<Type>asList(new Uint256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> name() {
        final Function function = new Function(FUNC_NAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> approve(String _spender, BigInteger _value) {
        final Function function = new Function(
                FUNC_APPROVE, 
                Arrays.<Type>asList(new Address(_spender),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> totalSupply() {
        final Function function = new Function(FUNC_TOTALSUPPLY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> allowTransfers() {
        final Function function = new Function(FUNC_ALLOWTRANSFERS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> transferFrom(String _from, String _to, BigInteger _value) {
        final Function function = new Function(
                FUNC_TRANSFERFROM, 
                Arrays.<Type>asList(new Address(_from),
                new Address(_to),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> limitedWallets(String param0) {
        final Function function = new Function(FUNC_LIMITEDWALLETS, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<BigInteger> decimals() {
        final Function function = new Function(FUNC_DECIMALS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> limitedWalletsManager() {
        final Function function = new Function(FUNC_LIMITEDWALLETSMANAGER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Boolean> issuanceFinished() {
        final Function function = new Function(FUNC_ISSUANCEFINISHED, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<String> manager() {
        final Function function = new Function(FUNC_MANAGER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> decreaseApproval(String _spender, BigInteger _subtractedValue) {
        final Function function = new Function(
                FUNC_DECREASEAPPROVAL, 
                Arrays.<Type>asList(new Address(_spender),
                new Uint256(_subtractedValue)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> balanceOf(String _owner) {
        final Function function = new Function(FUNC_BALANCEOF, 
                Arrays.<Type>asList(new Address(_owner)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> delLimitedWalletAddress(String _wallet) {
        final Function function = new Function(
                FUNC_DELLIMITEDWALLETADDRESS, 
                Arrays.<Type>asList(new Address(_wallet)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> symbol() {
        final Function function = new Function(FUNC_SYMBOL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> setLimitEnabled(Boolean _setLimitEnabled) {
        final Function function = new Function(
                FUNC_SETLIMITENABLED, 
                Arrays.<Type>asList(new Bool(_setLimitEnabled)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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

    public RemoteCall<TransactionReceipt> destroy(String _from, BigInteger _value) {
        final Function function = new Function(
                FUNC_DESTROY, 
                Arrays.<Type>asList(new Address(_from),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> isLimitedWalletAddress(String _wallet) {
        final Function function = new Function(FUNC_ISLIMITEDWALLETADDRESS, 
                Arrays.<Type>asList(new Address(_wallet)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> transfer(String _to, BigInteger _value) {
        final Function function = new Function(
                FUNC_TRANSFER, 
                Arrays.<Type>asList(new Address(_to),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setListener(String _listener) {
        final Function function = new Function(
                FUNC_SETLISTENER, 
                Arrays.<Type>asList(new Address(_listener)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> finishIssuance() {
        final Function function = new Function(
                FUNC_FINISHISSUANCE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> eventListener() {
        final Function function = new Function(FUNC_EVENTLISTENER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> increaseApproval(String _spender, BigInteger _addedValue) {
        final Function function = new Function(
                FUNC_INCREASEAPPROVAL, 
                Arrays.<Type>asList(new Address(_spender),
                new Uint256(_addedValue)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> isLimitEnabled() {
        final Function function = new Function(FUNC_ISLIMITENABLED, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<BigInteger> allowance(String _owner, String _spender) {
        final Function function = new Function(FUNC_ALLOWANCE, 
                Arrays.<Type>asList(new Address(_owner),
                new Address(_spender)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> setAllowTransfers(Boolean _allowTransfers) {
        final Function function = new Function(
                FUNC_SETALLOWTRANSFERS, 
                Arrays.<Type>asList(new Bool(_allowTransfers)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> ownerByAddress(String param0) {
        final Function function = new Function(FUNC_OWNERBYADDRESS, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> addLimitedWalletAddress(String _wallet) {
        final Function function = new Function(
                FUNC_ADDLIMITEDWALLETADDRESS, 
                Arrays.<Type>asList(new Address(_wallet)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setOwners(List<String> _owners) {
        final Function function = new Function(
                FUNC_SETOWNERS, 
                Arrays.<Type>asList(new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class))),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MCCToken> deploy(Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider, String _listener, List<String> _owners, String _manager) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Address(_listener),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class)),
                new Address(_manager)));
        return deployRemoteCall(MCCToken.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<MCCToken> deploy(Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider, String _listener, List<String> _owners, String _manager) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Address(_listener),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class)),
                new Address(_manager)));
        return deployRemoteCall(MCCToken.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCToken> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, String _listener, List<String> _owners, String _manager) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Address(_listener),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class)),
                new Address(_manager)));
        return deployRemoteCall(MCCToken.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<MCCToken> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, String _listener, List<String> _owners, String _manager) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Address(_listener),
                new DynamicArray<Address>(
                        org.web3j.abi.Utils.typeMap(_owners, Address.class)),
                new Address(_manager)));
        return deployRemoteCall(MCCToken.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    public List<BurnEventResponse> getBurnEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(BURN_EVENT, transactionReceipt);
        ArrayList<BurnEventResponse> responses = new ArrayList<BurnEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            BurnEventResponse typedResponse = new BurnEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.burner = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<BurnEventResponse> burnEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, BurnEventResponse>() {
            @Override
            public BurnEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(BURN_EVENT, log);
                BurnEventResponse typedResponse = new BurnEventResponse();
                typedResponse.log = log;
                typedResponse.burner = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<BurnEventResponse> burnEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(BURN_EVENT));
        return burnEventObservable(filter);
    }

    public List<AllowTransfersChangedEventResponse> getAllowTransfersChangedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(ALLOWTRANSFERSCHANGED_EVENT, transactionReceipt);
        ArrayList<AllowTransfersChangedEventResponse> responses = new ArrayList<AllowTransfersChangedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            AllowTransfersChangedEventResponse typedResponse = new AllowTransfersChangedEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._newState = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<AllowTransfersChangedEventResponse> allowTransfersChangedEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, AllowTransfersChangedEventResponse>() {
            @Override
            public AllowTransfersChangedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(ALLOWTRANSFERSCHANGED_EVENT, log);
                AllowTransfersChangedEventResponse typedResponse = new AllowTransfersChangedEventResponse();
                typedResponse.log = log;
                typedResponse._newState = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<AllowTransfersChangedEventResponse> allowTransfersChangedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ALLOWTRANSFERSCHANGED_EVENT));
        return allowTransfersChangedEventObservable(filter);
    }

    public List<IssueEventResponse> getIssueEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(ISSUE_EVENT, transactionReceipt);
        ArrayList<IssueEventResponse> responses = new ArrayList<IssueEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            IssueEventResponse typedResponse = new IssueEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._to = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<IssueEventResponse> issueEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, IssueEventResponse>() {
            @Override
            public IssueEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(ISSUE_EVENT, log);
                IssueEventResponse typedResponse = new IssueEventResponse();
                typedResponse.log = log;
                typedResponse._to = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<IssueEventResponse> issueEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ISSUE_EVENT));
        return issueEventObservable(filter);
    }

    public List<DestroyEventResponse> getDestroyEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(DESTROY_EVENT, transactionReceipt);
        ArrayList<DestroyEventResponse> responses = new ArrayList<DestroyEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            DestroyEventResponse typedResponse = new DestroyEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._from = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<DestroyEventResponse> destroyEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, DestroyEventResponse>() {
            @Override
            public DestroyEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(DESTROY_EVENT, log);
                DestroyEventResponse typedResponse = new DestroyEventResponse();
                typedResponse.log = log;
                typedResponse._from = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<DestroyEventResponse> destroyEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(DESTROY_EVENT));
        return destroyEventObservable(filter);
    }

    public List<IssuanceFinishedEventResponse> getIssuanceFinishedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(ISSUANCEFINISHED_EVENT, transactionReceipt);
        ArrayList<IssuanceFinishedEventResponse> responses = new ArrayList<IssuanceFinishedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            IssuanceFinishedEventResponse typedResponse = new IssuanceFinishedEventResponse();
            typedResponse.log = eventValues.getLog();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<IssuanceFinishedEventResponse> issuanceFinishedEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, IssuanceFinishedEventResponse>() {
            @Override
            public IssuanceFinishedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(ISSUANCEFINISHED_EVENT, log);
                IssuanceFinishedEventResponse typedResponse = new IssuanceFinishedEventResponse();
                typedResponse.log = log;
                return typedResponse;
            }
        });
    }

    public Observable<IssuanceFinishedEventResponse> issuanceFinishedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ISSUANCEFINISHED_EVENT));
        return issuanceFinishedEventObservable(filter);
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

    public List<TransferEventResponse> getTransferEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(TRANSFER_EVENT, transactionReceipt);
        ArrayList<TransferEventResponse> responses = new ArrayList<TransferEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            TransferEventResponse typedResponse = new TransferEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.from = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.to = (String) eventValues.getIndexedValues().get(1).getValue();
            typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<TransferEventResponse> transferEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, TransferEventResponse>() {
            @Override
            public TransferEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(TRANSFER_EVENT, log);
                TransferEventResponse typedResponse = new TransferEventResponse();
                typedResponse.log = log;
                typedResponse.from = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.to = (String) eventValues.getIndexedValues().get(1).getValue();
                typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<TransferEventResponse> transferEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(TRANSFER_EVENT));
        return transferEventObservable(filter);
    }

    public List<ApprovalEventResponse> getApprovalEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(APPROVAL_EVENT, transactionReceipt);
        ArrayList<ApprovalEventResponse> responses = new ArrayList<ApprovalEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            ApprovalEventResponse typedResponse = new ApprovalEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.owner = (String) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.spender = (String) eventValues.getIndexedValues().get(1).getValue();
            typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ApprovalEventResponse> approvalEventObservable(EthFilter filter) {
        return web3j.ethLogObservable(filter).map(new Func1<Log, ApprovalEventResponse>() {
            @Override
            public ApprovalEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(APPROVAL_EVENT, log);
                ApprovalEventResponse typedResponse = new ApprovalEventResponse();
                typedResponse.log = log;
                typedResponse.owner = (String) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.spender = (String) eventValues.getIndexedValues().get(1).getValue();
                typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ApprovalEventResponse> approvalEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(APPROVAL_EVENT));
        return approvalEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> issue(String _to, BigInteger _value) {
        final Function function = new Function(
                FUNC_ISSUE, 
                Arrays.<Type>asList(new Address(_to),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> burn(BigInteger _value) {
        final Function function = new Function(
                FUNC_BURN, 
                Arrays.<Type>asList(new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> burnFrom(String _from, BigInteger _value) {
        final Function function = new Function(
                FUNC_BURNFROM, 
                Arrays.<Type>asList(new Address(_from),
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    @Deprecated
    public static MCCTokenLocal load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCTokenLocal(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static MCCTokenLocal load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new MCCTokenLocal(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static MCCTokenLocal load(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        return new MCCTokenLocal(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static MCCTokenLocal load(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return new MCCTokenLocal(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    protected String getStaticDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

    public static String getPreviouslyDeployedAddress(String networkId) {
        return _addresses.get(networkId);
    }

}
