package blastertest

func (suite *blasterTestSuite) TestSingleBlaster() {
    // Act
    suite.setup(1, 10, 5)
    suite.start()
    suite.wait()

    // Assert
    suite.assertTotalRequestsWithin(0.80)
    suite.Greater(suite.successfulRequests(), 0)
}

