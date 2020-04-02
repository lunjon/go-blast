// +build many

/*
many_test.go -- test running many blasters, hence the build tag.
*/
package blastertest

func (suite *blasterTestSuite) TestManyBlastersHighRate() {
    // Arrange
    suite.setup(100, 10, 5)

    // Act
    suite.start()
    suite.wait()

    // Assert
    suite.assertTotalRequestsWithin(0.80)
    suite.Greater(suite.successfulRequests(), 0)
}

