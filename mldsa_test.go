package mldsa_test

import (
	"encoding/hex"
	"testing"

	mldsa "github.com/jasoncolburne/ml-dsa-go"
)

func testParamsRoundtrip(params mldsa.ParameterSet, t *testing.T, skLength, vkLength, sigLength int) {
	vk, sk, err := mldsa.KeyGen(params)
	if err != nil {
		t.Fatalf("error generating keypair: %v", err)
	}

	if len(sk) != skLength {
		t.Fatalf("unexpected sk length: %d != %d", len(sk), skLength)
	}

	if len(vk) != vkLength {
		t.Fatalf("unexpected sk length: %d != %d", len(vk), vkLength)
	}

	// for KAT testing
	// hexString := "9398dba05f34ea8e6bac58fd37f83d3e38ceb477f4653fe895a63be478c1dd8328ba5be45acb5cdcf1d4179bdb517442eb4a8974a363698e37e574818137ea4cd8984272fd41f0270bdf82c8ec4535078d70af1c66b3cfde7c7ea8d3fec252ff9f6bf39fdf78eb2d6f4c3bf71c8eeab654854ba0c917f29a8dafe2407eea9d986adf5393443806ed2f6f4d8e48139af331060f1361a535b634e600166417358c4867bf635c4f096cbba84004cf37bc1739e7cef67330a7769e8fece463c08f13946b65b46b042a06419041c63f17509d9811787511961d26463b1a453ca855a39a9191598a93675ac4ec8d82387793517ec37378c1e6f64c040c3b0df2c8ef960db8d0df64077644b39f5ca68b1eed61da95a464a7af30673ddaeec7a8573dcde37fb42d25fa0b9f7c6fed21b3be8a536a6982f1a988ff2b0f034adcf920bff805677e1f9452d7320e8095f22c6e89dd0c333ced96444baa9aa3fd1225733d8dd2ab9be7a95d6892b02818e1962037ed4a984bee33481c3044bb9d915444726196cafeab3289a2fa0ee908ea28f180cebd4f266a121ac74e2d1aae6438e26027313bc446857a56271b9d7e9db9e35964a0a37b088e2b2b67d3cbcf2e72a1a18caff3d3114f1abafc3375dfb7e9709a3bd3d6d77049c3b741b084316241d16d03f40a71f6714de91c33ca34aac05dcbf128f5c673b92c743b71060190e934b4896bb097047ebe74e08b9429823dc4e525f08bf3aede9f202e50846af0330bef16c412493b81d4d547229e73f03133ef723b06e406f8783dc8d47ae753ee1634c3eb801f9a5fb2a80f38e7fccf0f648ef91b84b62a8625af8d6bb98061ade0024ab90ed20efab6c3cd5f4848aa187eaf0b4450e9631c153fc83f62d63422b980bfae2eb0d39d9cdf8818138e7c21b064344a5ff844115de7bd817108f304842c3c166b843e0e144da444805350131dd9c9e08d2459605490859214f0d8225f8f9ea4b09fb90cd1ed1368b13288757c412913a3580caa12f8119ba802c857095fe5b1aec31d75f82af6769ac4fea18d4aa377b435b1b4e38be0419785ccfdd98f8af5c0f542bc1aa18608992f642e0fba633fce569e9d55c59999dd66cf57f239f273507043b3b647fe077b601874603585f9d607b6f8c88a1140be08ab588c67d65a4d9173444ba703ba22f3f3db24d6ccefbf48b3a44d179bbb992f4326d3fd24e3983ca6e2112f4b52af234718559032b3e3faa85feefa40270ddd21a160a424cba945c58aa2e30f42dcc7f86a50b2fd2dfd9d3c12168bc027611bd7edc6e6ce16f2ec4376ad1654bdcdcf432d5fd4f3ef52c84ea81c7f047fd833fa4a8b97fa2c519e066c933c8f89b922722757da301d4426f5ed8b86ea8d3102048f500f94f43f6a40bacf1cdbb5f19162c515aeff43acf1fb9519e13e0b468860ab5cf008e8acb3350430fc739f9a0c833d829052021f2d27d2c507bb65429230f85046f68c4f1920e1a85b90949edd3b34d5d4e622d7c36c9c41b88cf417ceed2970b0b4b4507ac757ce109931cd09b422473ceeb52fc33bce3af8a8b5aa9dcef27be840dd23a5cbfee02baffa9db5228de78f92046b71e465ab888a6877ad5f5d7059baf9e26ebabf126cfd9515486a0763add8f4938704b8d60ef39f89ace3730bdbe654659c26f38f225fb1385c786b35e359ad945301ba576bb870021dd5c07219afab2bcaa670991f461731e7aa2ff6ceb23f36d8ac4cda53631693ba8c7556264d7d53d35bcd79d74868e096d6ad5ad2d27cdd5b9e96f40d20f0616ea2facc8ad016477e5b652b5359e14d309ea7a984f9613b50f2180ceb66d686c29a900657a2"
	// expectedPk, _ := hex.DecodeString(hexString)
	// if subtle.ConstantTimeCompare(vk, expectedPk) != 1 {
	// 	if subtle.ConstantTimeCompare(vk[:32], expectedPk[:32]) != 1 {
	// 		fmt.Printf("rho mismatched:\n%s\n%s\n", hex.EncodeToString(vk[:32]), hex.EncodeToString(expectedPk[:32]))
	// 	}
	// 	offset := 32
	// 	delta := 320
	// 	limit := offset + delta
	// 	i := 0

	// 	for offset < len(vk) {
	// 		actual := vk[offset:limit]
	// 		expected := expectedPk[offset:limit]

	// 		if subtle.ConstantTimeCompare(actual, expected) != 1 {
	// 			fmt.Printf("polynomial %d mismatched:\n%s\n%s\n%v\n%v\n", i, hex.EncodeToString(actual), hex.EncodeToString(expected), actual, expected)
	// 		}

	// 		offset += delta
	// 		limit = offset + delta
	// 		i += 1
	// 	}

	// 	t.Fatalf("bad pk")
	// }

	// hexString = "9398dba05f34ea8e6bac58fd37f83d3e38ceb477f4653fe895a63be478c1dd8399fb4133d381432e87aa871bfb6f59fdb6c1aca1516a93264a7bfe510af1472da73e74ff36af20834e4e2e65e166a657db69cdb8414ef69059beff232b3e092f1cf9da3afce978ebd8605dc78756160ffe04b0b30ced19596c451cae5d017379cb166024490258404ce4304a03a8050499618924695328410ab301898288a1c4416248525c145014497100b7841c1570e3c4800238201b93648bb4081327241149888ac68ca3008de318484a3642e1b641e418608b282608809020252002151000272840446c2028065a14825a4081cb2028e3428d58002c22c988cc2205ca861113a9449cb024d0028cdc08469420514a8670d0382222a509dac02dc1286d1920091827301bb07020402840c48dc9c46843a8858b362500b291188661093888640690529065e1a805804661099665809828db346652126c8b08110ca8416302808ca8059a14209938260394049b180a18220404a9890bb45054226d138649140786e3208988223198442d40200509b6440b358849c02ddac24181224920a429e1468d00a9891030666002921bb24c239425101352a3486d1482000c406802068d1894644a488e99c42c00313008864854388923128863862180c20921458e23833014c28851463201031223030c48406a99c8490244211a9845528851cb482294a61049b491e4920c18b2518904229c4620d192859c405092204652b444518470daa68143c88cdc068ac0066d901828123681d0b280da348a59307192c64094b89020b64923312de0386e5cb0311aa6691a37459a388d24a920202100cc284a8c046600c48950c8259b36620401484a4428d286809844526112810c434ad93812cc080180448c14386900a60d5c126c50c86898a08ca2924d63922059042154846c5b2084e4465001446d82c80ccc20620b056ec2929103b969e2a84843486d189444222126da124cda266112b888a1c091d2069143369260c2110c0186503252dc80814382294b2849dab26910386c1a46849204040939720aa148248504843026e0c6290c1466c2c845819264c1804424226d22b46d1447011b3632c822064186241b45465c3841ccc809c24252c836080ba664428689e22642d1b02044060809b80ce044616128295b24300b17499094680a218ed136224a022c0a39859914880b1568e10051008529a18220da9028d092449a80250b106acaa2315b344abf6187bb73fedad39d21f321fe92ad2e999de228e6ae22b33bef9b28554c1d38a7564ae71f498be650a94fe4850a5c998a62e19c803ef34b9105b815c83de1743a1486581f56461173291ad807b0e4b4202e3206a343e4f88c426a086a74f7582a72db25e4629305523402804ce2662da2e0b328d3eaac6069a4a2909c8468c95a563222dd792b9a49ebec62366cece6ddc9dbdda03ee1ce569cdcf020a4d62a10b89470ebe31b1e01ed320e2b6eb46ed17c1b0159a1b6e83f557433a564428c44605e4860de56b8efac04a17d7673c5c79e248bcfd5e209790788371b21d665567d191d0b1f41d70168e411b6b33367af715216cce3210697e18df66e3490bd0d63728fd47792c63385145f81cad53565630cd54ed621f3fab3e29972a3ac64ef1576639bc57114a11fd26fa764fad7563b7a08fa24a4ba6607056471bda358d8c5babdc3bd6b10b0c121e1ab1dd944cc290d4fdd810e160376cf640d9365e14f3ed9d77342a9336cb7b0f0b08c9e7fc9ba024d02d1e38f75f6df0e4876bf77e64b12965af1b5e7d081557aec33af07d391ac42a729879062eafd3e0d8375cfb4cf8c69359cbac557093f5fb5cc1f7957b44c995c595d5eede3045d462f37530ff0701a77817c49abcce9af92a628d74d973d84aa182057288d862b2ea007ee7d5e61994db47a827fee93dadb494e72208f949422656fe392f6a17eca044f9cb44f5a96a3cbfcf97fe365f06f2ee54f3d5ab0780a057fcbc2c7ea4c9448c852331dcaf45f33e2feaf10f69cdf79ab2146fda7f053ffd6bbc33f99a0440d1c58d6a47be5cf6ba497f71d98a97993ebb7b3880bae2f67fb70ad957bf6ca0211aa5d696b7b5c118d9110a4c68c4e232ec296d3a42cbc54260b0b27ea1b977468bb3fe5fda6d545809cc56af7edc5a70ca1e4d9f1e330fc17ff10a4cd18d33e416c690c4c4b93735653c92bfd5da9eea1b75cc99119f1f6d89fdc312570a76eefa7490ce5b64b5d989e6b5ee051ac870783ca33e535a361f10f530b631824006bb622f43f839d103a9466a011de9a704889d964b118da368581cec167d03ca2f7ea26e5174e925a0aa0a0f4f5d16f1ba3c10354d3a6cd4a3fedac1ee81391d6c020a3932d5ce0433958d1200c02424ffb7f3001a8677f17c278f71e17317cce993632953b91e2855c24a9eaf92231914b88900c0b08b44cf03851242e21b4895f16968d9ced49748932b3a21ebd426be3945e3e53fca08bba7eb36b7e2f7b1cf5dfd1deb9399c14110544a1e4b565ca7b87181e60944c263d623137f8f8e0f0001284b7530fa430081dd19dde692a6b3ba7af6405a95b9a5adffb6e3502bb35be8a7e524a412f735b13cac770d54e1b5e232f999098c02932d88303d141bf044b238673c154152821ce21c69aa3a8abf3ae7b424cb29b5c7bd895ed16888e1e784b25e374fca856e9758a33d41a2ddb23a71da7d6803a62abc574b6f62c0dbe73605770c8c248c0caa23d11e4cc0b0fed844d905afc84aa2429faf58ff5b08b11b5816a045be1109ec6e55b150b5e5690bb628f3548a3d681e053d58d18a8b8d1d83fa0c1c9ca710d0165fa696e9390d47ba768ade0b8ce7e46fbf5484d1ecd011944c0e6ed474cfc50eeca5dbeba4b0399a6a035e85be2792885435acb280f8b620373a2a0e9a1f9d6d903b5b67dda376993c69cd6621fa38c591ce98b9e30b8ea22763af2d7e6735dfb98087eb1c49053f31de79f3a6f6f2224899e0d3ea617b97b797f54bc0b17a938ea4fc1d89bfcd8a4523f57bc03fb391cea2ea8294c86d575664f55e00ba904a947bcfc2b2676c4905a49f7ccb5c1d914e8c15bb2a02cd73662a0e7bd2b6fb016a67db3679b0cb40f2bb9f344fc1f59eb3e84e9247ffbb4961dea61e5c614984c0f89d37e9dd68e180381e0bdc19f2c827b908348252a318851090359e5bd7a12c4b25730973c5423188bda92b017d54ce93e757a0812264ebc3f1addde07957be0e19154c5ce827e4e97de7c66d71ec7eca32c78bfe0eb116efc269cb4fb872159758d02811ccfb043b057d6f0ea08e1f4afa8b2afb98449a004488b14a80835b4be74acecb484133379f6cd22aad087c221f7109c6907ce1063af687f9ce1ab6a87b41a71dc5f32b825aa47fe512a49f0ccdbfaa1e7501d25f9bbac6ecddd904def9411f5d334479f0167bf53ed4e161a7c5fbc8f4c80da406c17ba6865db2433abc6a5163a5f2df6c25db3e4c56392215b477114bab2df031e1a4fdaf98a79af8a6de5bf1e39e3b7f26dd3d3370e5a9d98798c109414ac14279373e4f10e007234c7a0202fedd22fe45faf2f4b734dcd7bb64a0d3eae"
	// expectedSk, _ := hex.DecodeString(hexString)

	// if subtle.ConstantTimeCompare(sk, expectedSk) != 1 {
	// 	t.Fatalf("bad sk")
	// }

	message, _ := hex.DecodeString("89b0c4b23019af3498a27da290892d981dd59fa08993bc05da21e1d72503664c98cadefc061d176d0b44bcab049bb540e0680a58bdad0d16316f772d44d47281")
	ctx, _ := hex.DecodeString("09764e76473cc969442691dd0574afdd")
	sig, err := mldsa.Sign(params, sk, message, ctx)
	if err != nil {
		t.Fatalf("error signing: %v", err)
	}

	if len(sig) != sigLength {
		t.Fatalf("unexpected sig length: %d != %d", len(sig), sigLength)
	}

	// expected, _ := hex.DecodeString("f7b2093bd16daf10a2dd31c3b6a23ca4ede151b5ef75f8d5f70f299f21df188ce4840ed664efb399978d16cb3f627a25ec9c65711560744566164cf10880ff27af6fe906d471529dc196c21ae50417961cbabb8898aabf1e212585727540a0acc99d2bd7638f68d7178ad2f93bab0718e1442eca80af97d7fa5d1f8f5855565c7c96eb4c6c60a718931857a6875ffba4c6ade6b1c7d2a943621d3b3fdb1ae2fd79c92d5470893241da3a895e90dd37cf3114e827209f017251272a115fff0178f5f44434de39a77f3da0dddeccf205476a2b47b48350e30ce8797e34c5e17bb39b08a25394823c7e93f821536b8a734c6dcea5bb7c53fc892f9e76e53e6cefbc3faa4691da5e1b9a3486f79f9a1b39eb020ee9358d4953b78b1a08625848307127a395fb949ad4544ea739a25c5b3e140120d80c8ae3c7e15900afac4d1e7028900039b15afb8088112ec4150176d34cbdc96ab2f3b27cfe63a75638dca3ab926c50f68aca7fcc2de87413280ee89795dafffce39b78a08a1d1a5c60cd8b617d3db517d3144880c3cca7bf3262d70acf5304f7eb706e16ef45452cc786da267d0091ffb45812f12858203b5f6deabc3f21bb5d20f6ad1f0bb30f48e26090103c9fdb97502de8f94464e9e8ac02e76252f3b793068aa5880718c04a89aebdb61dfeca55bf054e1dfa97d4f5c7f9d06d10256e8f5ac3e1e5d6852e108d7e74126391081bb43868a8689a87ab6e3ba2e94b49795fc4be455c2aa54de52c97eb0099f3e3b0ec22e56a0b5ff281a520d50bdd006d0f591687d05d87435a95c301b895a352b473b6dd8941e89388443c47bf4fd7d4a80d16cf7f3eb354e77fdf184bc6262c2ac175e26f3777470962051befc2020bed0093145d9544905d63624ab9ccf0e3b5a1305c4d2b39d59fd49bf97955544c9126611e552befb7e0464b22e78acbcfa88f8aa3e2e8d3a5bada0196d1d125f44fe785829e5a193176453569f9d07813d30eb81578a5b8d8f32eafd2032860dc68555db62a3c84b8a8106fa451ae5ebfbd406257003717b57f970f2e2ca8190897e34741c6c2dc472a6fe79ded9f8edbf64d431ed8db5c8ea76506a2610af01a9aafcaf007c39dacb3e55db9adf49320366237f33de768669327c2c2c3ee5678a8484ad11ad73567b72f8e14432a58e261c8ce0656922961c99cbd1674ea8fa464697756840194dc420b3d74301fcbdebc8bb562db5930f7c5eabdb452240f92a68972daafb4d002f268ef52abbc2e15612d47b9453b27d60bc735ebffadff9a6af94b495bc1873311a0e4b2e5a179a092d17ed38331dd17c85481c23830c5f52fd51880816672eec717d27485846763f1ebc5f66c0412a53bd35c4f43c23341fd4cc701c5b069bf99a6b7ea50a8a36386120b0d1fbacd07be6eb81ceb1d5e68677c7017f539c6a94289030ef5d1a59b43c239e9cce47c99dd6692d7c09f03a5f946d91728d00f04a3f71ddec80be8816cecb5750e61da830b2bf3d42f52f3d7f212ec76ea91d118060d0a8e3af8faca4f8f401f4238e85118dc7a31834188ef4869ad44506bf7ddd502ff7b6f32b7e1f8be9a47787c0f03b5bcb92a8a8b0bb61bdb8de3467c134ab0c320e00c0233bbd58301f1d3ed9c4d25e20f69a6a3e267397c771923f52c8f3bc152022083c9eb3dc60c0676e284aef097114ef8be7a7265d37a38de3e1a019179bfdf3337914e580bfe95573eb7d6e13fee33a71004584083a4536685c2bdb23c6836db4c071734f9aaba3793bbb66997c39554c6fbb4b88b7129c76a808a7e67e3aec052826e14bfb6d67da84bb8e2046460fc94606c94295703ba81182c962ff008a028e97ed12c511bcd3d834d594f7239940baad3b229d611976407839051e8feba8cab2046234dd067fc0d13c2e41f300a92d575c41f64e1a6d144f32585be1a44f941730f6dbff3e5fb580ad9d2bd8347f7f69714a0cec35764186648ff8633ee03ca8f32b0f0f9359f4dca219c27b273a7a7d067f115673e48819ca8dfba6d53bb632394cddf626fb4031316942905a7545da722f25023a7fa11ddf27f92956156b8cc21f50262af4d4f3fc6ce8a3b9805b91b866fdbaccca6d2b4c5a474c4c7d8c1a1eb5e67091f4f98bb60aa83d68e583a6d4f1c3c8b6506cd5af4d55318761bbceb20138d7e7eeb50fc8dc507acae3eb5e13ff920abec78d2f600dc2ce56fd0f43d3c6f50f8cbcdbaf28597d5ab910a477049ca30fd4135fd35b7e86dc421c4e5963099ef19251c9371ba0a436fa681ba72fed6b9f170372402011976711b56f9fbf73a7aaa621c3eb4e3d1ec87809e02bc783a10541bf964830c19d487251ca789947dec292a28cbb30cd374728efa96d0ec4c90f54e2745f5b5818e0893fc69bcc41be2903aa10c53a4b12f48a031373a1c7c19f329ea3e1756378bead6a711e0f06b2b8c7cceba58acc9c644950d47b6b98cc5415828fb9c281b39e22bc610a2a3d33efc73a93b16e327e52ef23ebf3097a2699a86c951a9d8b4d25375534776bdf9d7cb3eefdffa90c3376cde30f508a4534c7faf2f71dca23f13ddb54e31fd4e588a8be30629083eaa3c368c44f68fccc5d5c15493fd2b9f9c8499631df2e13ca27b2a2fe3c870c49c763888599a57061aab9a0ad2686044287ab100a36cb336a5daeb1dea8ac1f4f87c2cc8af00d863cbc33fdfaf0e4a069bb7d3b74d791acac185262d1598c52645cf174a97c384ef6d13e0938015133601e2f98b41541f5312806fd9190d6a6f0d64b2097cb21dc8a6cb6b9de9693c9b4559b96b0c83a44e052365979451c2c81d2dc462b31f167fbc9ee9d941e142760e2e11a63251f682269a245e27184608120ece898dd68a7fe687323c461fa212f8d229674015d200ae22caaf51561d58b2292f11f406cd5205740be8cc5d49e2cd5af549627e858a2a6755ab8638dd505216e2d4bd33e5f7d086001887e80da9e58df57e16b4d699f0b5ff314bee300f97a51df35a9afb5b509bd8d51afadfa6dc7c98cfde7143542adb0bcd9c304ebb5d63ef00fe1265cfcb0f4c1a10e45465249387c67532e77fd3eba03bab6d442fca2d5790fbb039e0db941fd9818f2a428b8b20671ca88ad8cebca06cf95c8d9e89852d96a75565b97f6daa218507b829eaad4a62da35f16710319629ffb3ee37c1e3126dc33112edcadd0a8e47aa3f3d1898f94ea8e8d5d6dccd0d0f8f81e9b2392af9dc4f9aa77039cd205a7a6380ad65bcf19fa89044d23397585aefcaaa3702a92b55e82febc7e7553be5e75243688d9d8e2b7b060814162a2b323340474b4f506c7b7c8e97a2abb1b2bbe7f303050d282a2b3b43484e7d7eb5d8f9fc323f4c5559677273aebcd3e5ef060a1117475455575d656e70797f8c9094a2a5c1c5d4ebf800001929364e89b0c4b23019af3498a27da290892d981dd59fa08993bc05da21e1d72503664c98cadefc061d176d0b44bcab049bb540e0680a58bdad0d16316f772d44d47281")
	// actual := make([]byte, len(sig))
	// copy(actual, sig)
	// actual = append(actual, message...)

	// if subtle.ConstantTimeCompare(actual, expected) != 1 {
	// 	a := actual[:32]
	// 	e := expected[:32]
	// 	if subtle.ConstantTimeCompare(a, e) != 1 {
	// 		fmt.Printf("rho mismatched")
	// 	}

	// 	t.Fatalf("bad sig")
	// }

	valid, err := mldsa.Verify(params, vk, message, sig, ctx)
	if err != nil {
		t.Fatalf("error verifying: %v", err)
	}

	if !valid {
		t.Fatalf("signature not valid!")
	}
}

// func calculateEntropy(data []byte) float64 {
// 	frequency := make(map[byte]int)
// 	for _, b := range data {
// 		frequency[b]++
// 	}

// 	var entropy float64
// 	for _, freq := range frequency {
// 		p := float64(freq) / float64(len(data))
// 		entropy -= p * math.Log2(p)
// 	}
// 	return entropy
// }

func TestMLDSA44RoundTrip(t *testing.T) {
	testParamsRoundtrip(mldsa.ML_DSA_44_Parameters, t, 2560, 1312, 2420)
}

// func TestMLDSA65RoundTrip(t *testing.T) {
// 	testParamsRoundtrip(mldsa.ML_DSA_65_Parameters, t, 4032, 1952, 3309)
// }

// func TestMLDSA87RoundTrip(t *testing.T) {
// 	testParamsRoundtrip(mldsa.ML_DSA_87_Parameters, t, 4896, 2592, 4627)
// }
