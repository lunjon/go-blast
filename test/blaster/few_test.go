/*
few_test.go -- run a few blasters, no more than ten.
*/
package blastertest

func (suite *blasterTestSuite) TestFiveBlasters() {
    // Act
    suite.setup(5, 10, 5)
    suite.start()
    suite.wait()

    // Assert
    suite.assertTotalRequestsWithin(0.80)
    suite.Greater(suite.successfulRequests(), 0)
}

func (suite *blasterTestSuite) TestFiveBlastersHighRate() {
    // Act
    suite.setup(5, 10, 5)
    suite.start()
    suite.wait()

    // Assert
    suite.assertTotalRequestsWithin(0.80)
    suite.Greater(suite.successfulRequests(), 0)
}
