Feature: Terraform Host Provider works
  As a user of Docker and Uberstack, I need to be able to create and destroy hosts.

  Scenario: I can create a new host
    Given UBER_HOME is set
    When I execute 'uberstack host up terraform-host01'
    Then the output contains 'created.'

#  Scenario: I can run a single command on a remote host
#    Given a running host 'terraform-host01'
#      And a known SSH key
#    When I execute 'uname' via ssh on host 'terraform-host01'
#    Then the output contains 'Linux'

#  Scenario: I can run a single command on a remote host
#    Given a running host 'terraform-host01'
#      And a known SSH key
#    When I execute 'echo ABC > foo; cat foo' via ssh on host 'terraform-host01'
#    Then the output contains 'ABC'

 # Scenario: I can delete a host
 #   Given UBER_HOME is set
 #   When I execute 'uberstack host rm terraform-host01'
 #   Then the output contains 'deleted.'
