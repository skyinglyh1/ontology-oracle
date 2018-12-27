from boa.interop.System.App import RegisterAppCall
from boa.interop.System.ExecutionEngine import GetExecutingScriptHash
from boa.interop.System.Runtime import Notify, Serialize, Deserialize
from boa.builtins import *
from ontology.interop.Ontology.Runtime import Base58ToAddress

oracleContract = RegisterAppCall('d1ea6327aa344f37e5e3c7ede3125c09a6cb4b45', 'operation', 'args')

def main(operation, args):
    if operation == 'genRandom':
        return genRandom()
    if operation == 'getRandom':
        if len(args) == 1:
            return getRandom(args[0])
    return False


def genRandom():

    req = """{
		"scheduler":{
			"type": "runAfter",
			"params": "2018-06-15 08:37:18"
		},
		"tasks":[
			{
			    "type": "httpGet",
			    "params": {
				    "url": "https://data.nba.net/prod/v2/20181129/scoreboard.json"
			    }
			},
			{
				"type": "jsonParse",
				"params":
				{
					"data":
					[
						{
							"type": "Array",
							"path": ["games"],
							"sub_type":
							[
								{
									"type": "Struct",
									"sub_type":
									[
										{
											"type": "String",
											"path": ["gameId"]
										},
										{
											"type": "String",
											"path": ["vTeam", "teamId"]
										},
										{
											"type": "String",
											"path": ["vTeam", "score"]
										},
										{
											"type": "String",
											"path": ["hTeam", "teamId"]
										},
										{
											"type": "String",
											"path": ["hTeam", "score"]
										}
									]
								}
							]
						}
					]
				}
			}
		]
	}"""

    res = oracleContract('CreateOracleRequest',[req,Base58ToAddress('AKQj23VWXpdRou1FCUX3g8XBcSxfBLiezp')])

    return True

def getRandom(txHash):
    res = oracleContract('GetOracleOutcome', [txHash])
    if not res:
        return ''
    a = Deserialize(res)
    b = Deserialize(a[0])
    Notify(b)
    Notify(b[0])
    Notify(b[0][0])
    Notify(b[0][0][0])
    return True