package dht

import (
	"encoding/hex"
	"fmt"
	"testing"
	"math/big"
)

// test cases can be run by calling e.g. go test -test.run TestRingSetup
// go run test will run all tests

/*
 * Example of expected output of calling printRing().
 *
 * f38f3b2dcc69a2093f258e31902e40ad33148385 1390478919082870357587897783216576852537917080453
 * 10dc86630d9277a20e5f6176ff0786f66e781d97 96261723029167816257529941937491552490862681495
 * 35f2749bbe6fd0221a97ecf0df648bc8355c7a0e 307983449213776748112297858267528664243962149390
 * 3cb3aaec484f62c04dbab1512409b51887b28272 346546169131330955640073427806530491225644106354
 * 624778a652b23ebeb2ce133277ee8812fff87992 561074958520938864836545731942448707916353010066
 * a5a5dcfbd8c15e495242c4d7fe680fe986562ce2 945682350545587431465494866472073397640858316002
 * b94a0c51288cdaaa00cd5609faa2189f56251984 1057814620711304956240501530938795222302424635780
 * d8b6ac320d92fe71551bed2f702ba6ef2907283e 1237215742469423719453176640534983456657032816702
 * ee33f5aaf7cf6a7168a0f3a4449c19c9b4d1e399 1359898542148650805696846077009990511357036979097
 */
func TestRingSetup(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id8 := "08"

	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(&id0, "localhost", "1111")
	node2 := makeDHTNode(&id1, "localhost", "1112")
	node3 := makeDHTNode(&id2, "localhost", "1113")
	node4 := makeDHTNode(&id3, "localhost", "1114")
	node5 := makeDHTNode(&id4, "localhost", "1115")
	node6 := makeDHTNode(&id5, "localhost", "1116")
	node7 := makeDHTNode(&id6, "localhost", "1117")
	node8 := makeDHTNode(&id7, "localhost", "1118")
	node9 := makeDHTNode(&id8, "localhost", "1119")

	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)


	for i := 0; i < 10; i++ {
		node1.stabilize()
		node2.stabilize()
		node3.stabilize()
		node4.stabilize()
		node5.stabilize()
		node6.stabilize()
		node7.stabilize()
		node8.stabilize()
		node9.stabilize()
	}
	for i := 0; i < 10; i++ {
		node1.fixFingers()
		node2.fixFingers()
		node3.fixFingers()
		node4.fixFingers()
		node5.fixFingers()
		node6.fixFingers()
		node7.fixFingers()
		node8.fixFingers()
		node9.fixFingers()
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

}

/*
 * Example of expected output.
 *
 * str=hello students!
 * hashKey=cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 */
func TestLookup(t *testing.T) {
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node2.addToRing(node1)
	node3.addToRing(node1)
	node4.addToRing(node1)
	node5.addToRing(node1)
	node6.addToRing(node1)
	node7.addToRing(node1)
	node8.addToRing(node1)
	node9.addToRing(node1)

	for i := 0; i < 10; i++ {
		node1.stabilize()
		node2.stabilize()
		node3.stabilize()
		node4.stabilize()
		node5.stabilize()
		node6.stabilize()
		node7.stabilize()
		node8.stabilize()
		node9.stabilize()
	}
	for i := 0; i < 10; i++ {
		node1.fixFingers()
		node2.fixFingers()
		node3.fixFingers()
		node4.fixFingers()
		node5.fixFingers()
		node6.fixFingers()
		node7.fixFingers()
		node8.fixFingers()
		node9.fixFingers()
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	str := "hello students!"
	hashKey := sha1hash(str)
	fmt.Println("str=" + str)
	fmt.Println("hashKey=" + hashKey)

	fmt.Println("node 1: " + hex.EncodeToString(node1.lookup(hashKey).nodeId) + " is respoinsible for " + hashKey)
	fmt.Println("node 5: " + hex.EncodeToString(node5.lookup(hashKey).nodeId) + " is respoinsible for " + hashKey)

	fmt.Println("------------------------------------------------------------------------------------------------")


}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            1
 * m            3
 * 2^(k-1)      1
 * (n+2^(k-1))  1
 * 2^m          8
 * result       1
 * result (hex) 01
 * successor    01
 * distance     1
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            2
 * m            3
 * 2^(k-1)      2
 * (n+2^(k-1))  2
 * 2^m          8
 * result       2
 * result (hex) 02
 * successor    02
 * distance     2
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            3
 * m            3
 * 2^(k-1)      4
 * (n+2^(k-1))  4
 * 2^m          8
 * result       4
 * result (hex) 04
 * successor    04
 * distance     4
 */
func TestFinger3bits(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0 := makeDHTNode(&id0, "localhost", "1111")
	node1 := makeDHTNode(&id1, "localhost", "1112")
	node2 := makeDHTNode(&id2, "localhost", "1113")
	node3 := makeDHTNode(&id3, "localhost", "1114")
	node4 := makeDHTNode(&id4, "localhost", "1115")
	node5 := makeDHTNode(&id5, "localhost", "1116")
	node6 := makeDHTNode(&id6, "localhost", "1117")
	node7 := makeDHTNode(&id7, "localhost", "1118")

	node1.addToRing(node0)
	node2.addToRing(node0)
	node3.addToRing(node0)
	node4.addToRing(node0)
	node5.addToRing(node0)
	node6.addToRing(node0)
	node7.addToRing(node0)

	for i := 0; i < 10; i++ {

		node0.stabilize()
		node1.stabilize()
		node2.stabilize()
		node3.stabilize()
		node4.stabilize()
		node5.stabilize()
		node6.stabilize()
		node7.stabilize()
	}
	for i := 0; i < 10; i++ {
		node0.fixFingers()
		node1.fixFingers()
		node2.fixFingers()
		node3.fixFingers()
		node4.fixFingers()
		node5.fixFingers()
		node6.fixFingers()
		node7.fixFingers()
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()

	fmt.Println("------------------------------------------------------------------------------------------------")

	/* 
	 * n            0
	 * k            1
	 * m            3
	 * 2^(k-1)      1
	 * (n+2^(k-1))  1
	 * 2^m          8
	 * result       1
	 * result (hex) 01
	 * successor    01
	 * distance     1
	 */
	k := 1
	m := 3
	a, b, c, r, d, rh, s := node0.testCalcFingers(k, m)
	if 	a.Cmp(big.NewInt(1)) != 0 || // 2^(k-1)  
		b.Cmp(big.NewInt(1)) != 0 || // (n+2^(k-1)) 
		c.Cmp(big.NewInt(8)) != 0 || // 2^m 
		r.Cmp(big.NewInt(1)) != 0 || // result 
		rh != "01" || // result (hex)
		s != "01" || // successor
		d.Cmp(big.NewInt(1)) != 0 { // distance

		t.Errorf("Finger k=%d failed (a=%s, b=%s, c=%s, r=%s, rh=%s, s=%s, d=%s)",k,a.String(),b.String(),c.String(),r.String(),rh,s,d.String())
	} 

	/* 
	 * n            0
	 * k            2
	 * m            3
	 * 2^(k-1)      2
	 * (n+2^(k-1))  2
	 * 2^m          8
	 * result       2
	 * result (hex) 02
	 * successor    02
	 * distance     2
	 */	
	k = 2
	m = 3
	a, b, c, r, d, rh, s = node0.testCalcFingers(k, m)
	if 	a.Cmp(big.NewInt(2)) != 0 || // 2^(k-1)  
		b.Cmp(big.NewInt(2)) != 0 || // (n+2^(k-1)) 
		c.Cmp(big.NewInt(8)) != 0 || // 2^m 
		r.Cmp(big.NewInt(2)) != 0 || // result 
		rh != "02" || // result (hex)
		s != "02" || // successor
		d.Cmp(big.NewInt(2)) != 0 { // distance

		t.Errorf("Finger k=%d failed (a=%s, b=%s, c=%s, r=%s, rh=%s, s=%s, d=%s)",k,a.String(),b.String(),c.String(),r.String(),rh,s,d.String())
	} 

	/* 
	 * n            0
	 * k            3
	 * m            3
	 * 2^(k-1)      4
	 * (n+2^(k-1))  4
	 * 2^m          8
	 * result       4
	 * result (hex) 04
	 * successor    04
	 * distance     4
	 */	
	k = 3
	m = 3
	a, b, c, r, d, rh, s = node0.testCalcFingers(k, m)
	if 	a.Cmp(big.NewInt(4)) != 0 && // 2^(k-1)  
		b.Cmp(big.NewInt(4)) != 0 && // (n+2^(k-1)) 
		c.Cmp(big.NewInt(8)) != 0 && // 2^m 
		r.Cmp(big.NewInt(4)) != 0 && // result 
		rh != "04" && // result (hex)
		s != "04" && // successor
		d.Cmp(big.NewInt(4)) != 0 { // distance

		t.Errorf("Finger k=%d failed (a=%s, b=%s, c=%s, r=%s, rh=%s, s=%s, d=%s)",k,a.String(),b.String(),c.String(),r.String(),rh,s,d.String())
	} 

}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            0
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0

 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            1
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            80
 * m            160
 * 2^(k-1)      604462909807314587353088
 * (n+2^(k-1))  682874255151879437996523461382311327142222978674
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996523461382311327142222978674
 * finger (hex) 779d240121ed6d5e8bd14b6529b08e5c617b5e72
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            120
 * m            90
 * 2^(k-1)      664613997892457936451903530140172288
 * (n+2^(k-1))  682874255152544051994415314855853423357775797874
 * 2^m          1237940039285380274899124224
 * finger       1180872106465109536036052594
 * finger (hex) 03d0cb6529b08e5c617b5e72
 * successor    f880fb198b7059ae92a69968727d84da9c94dd15
 * distance     877444087302148207702277795
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            160
 * m            160
 * 2^(k-1)      730750818665451459101842416358141509827966271488
 * (n+2^(k-1))  1413625073817330897098365273277543029655601897074
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       1413625073817330897098365273277543029655601897074
 * finger (hex) f79d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * successor    d0a43af3a433353909e09739b964e64c107e5e92
 * distance     508258282811496687056817668076520806659544776736
 */
func TestFinger160bits(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node2.addToRing(node1)
	node3.addToRing(node1)
	node4.addToRing(node1)
	node5.addToRing(node1)
	node6.addToRing(node1)
	node7.addToRing(node1)
	node8.addToRing(node1)
	node9.addToRing(node1)

	for i := 0; i < 20; i++ {
		node1.stabilize()
		node2.stabilize()
		node3.stabilize()
		node4.stabilize()
		node5.stabilize()
		node6.stabilize()
		node7.stabilize()
		node8.stabilize()
		node9.stabilize()
	}
	for i := 0; i < 500; i++ {
		node1.fixFingers()
		node2.fixFingers()
		node3.fixFingers()
		node4.fixFingers()
		node5.fixFingers()
		node6.fixFingers()
		node7.fixFingers()
		node8.fixFingers()
		node9.fixFingers()
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	node3.testCalcFingers(0, 160)
	fmt.Println("")
	node3.testCalcFingers(1, 160)
	fmt.Println("")
	node3.testCalcFingers(80, 160)
	fmt.Println("")
	node3.testCalcFingers(120, 90)
	fmt.Println("")
	node3.testCalcFingers(160, 160)
	fmt.Println("")

	/*
	 * n            682874255151879437996522856919401519827635625586
	 * k            0
	 * m            160
	 * 2^(k-1)      1
	 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
	 * 2^m          1461501637330902918203684832716283019655932542976
	 * finger       682874255151879437996522856919401519827635625587
	 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
	 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
	 * distance     0
	*/

	/*
	 * n            682874255151879437996522856919401519827635625586
	 * k            1
	 * m            160
	 * 2^(k-1)      1
	 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
	 * 2^m          1461501637330902918203684832716283019655932542976
	 * finger       682874255151879437996522856919401519827635625587
	 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
	 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
	 * distance     0
	*/
	/*
	 * n            682874255151879437996522856919401519827635625586
	 * k            80
	 * m            160
	 * 2^(k-1)      604462909807314587353088
	 * (n+2^(k-1))  682874255151879437996523461382311327142222978674
	 * 2^m          1461501637330902918203684832716283019655932542976
	 * finger       682874255151879437996523461382311327142222978674
	 * finger (hex) 779d240121ed6d5e8bd14b6529b08e5c617b5e72
	 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
	 * distance     0
	*/
	/*
	 * n            682874255151879437996522856919401519827635625586
	 * k            120
	 * m            90
	 * 2^(k-1)      664613997892457936451903530140172288
	 * (n+2^(k-1))  682874255152544051994415314855853423357775797874
	 * 2^m          1237940039285380274899124224
	 * finger       1180872106465109536036052594
	 * finger (hex) 03d0cb6529b08e5c617b5e72
	 * successor    f880fb198b7059ae92a69968727d84da9c94dd15
	 * distance     877444087302148207702277795
	*/
	/*
	 * n            682874255151879437996522856919401519827635625586
	 * k            160
	 * m            160
	 * 2^(k-1)      730750818665451459101842416358141509827966271488
	 * (n+2^(k-1))  1413625073817330897098365273277543029655601897074
	 * 2^m          1461501637330902918203684832716283019655932542976
	 * finger       1413625073817330897098365273277543029655601897074
	 * finger (hex) f79d240121ed6d5e8bd0cb6529b08e5c617b5e72
	 * successor    d0a43af3a433353909e09739b964e64c107e5e92
	 * distance     508258282811496687056817668076520806659544776736
	 */	
}
